package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	siliconFlowURL = "https://api.siliconflow.cn/v1/embeddings"
	embeddingModel = "Qwen/Qwen3-Embedding-0.6B"
	embeddingDim   = 1024
	batchSize      = 20
	chunkSize      = 500
	chunkOverlap   = 50
	minChunkLen    = 10
)

type embedReq struct {
	Input interface{} `json:"input"`
	Model string      `json:"model"`
}
type embedResp struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

type chunk struct {
	content  string
	source   string
	category string
}

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// Connect Milvus
	milvusAddr := os.Getenv("MILVUS_ADDR")
	if milvusAddr == "" { milvusAddr = "localhost:19530" }
	c, err := client.NewClient(ctx, client.Config{Address: milvusAddr})
	if err != nil { log.Fatal("Milvus connect error:", err) }
	defer c.Close()

	// Drop and recreate collection
	has, _ := c.HasCollection(ctx, "knowledge_chunks")
	if has {
		c.DropCollection(ctx, "knowledge_chunks")
		fmt.Println("Dropped old collection")
	}

	schema := entity.NewSchema().
		WithName("knowledge_chunks").
		WithDescription("RAG knowledge chunks").
		WithField(entity.NewField().WithName("id").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
		WithField(entity.NewField().WithName("content").WithDataType(entity.FieldTypeVarChar).WithMaxLength(4096)).
		WithField(entity.NewField().WithName("source").WithDataType(entity.FieldTypeVarChar).WithMaxLength(64)).
		WithField(entity.NewField().WithName("ref_id").WithDataType(entity.FieldTypeInt64)).
		WithField(entity.NewField().WithName("category").WithDataType(entity.FieldTypeVarChar).WithMaxLength(128)).
		WithField(entity.NewField().WithName("vector").WithDataType(entity.FieldTypeFloatVector).WithDim(embeddingDim))

	if err := c.CreateCollection(ctx, schema, 2); err != nil { log.Fatal("create collection:", err) }
	idx, _ := entity.NewIndexIvfFlat(entity.L2, 128)
	c.CreateIndex(ctx, "knowledge_chunks", "vector", idx, false)
	c.LoadCollection(ctx, "knowledge_chunks", false)
	fmt.Println("Collection created and loaded")

	// Read and chunk all policy files
	policyFiles := map[string]string{
		"return_policy":   "return_policy.md",
		"shipping_policy": "shipping_policy.md",
		"payment_policy":  "payment_policy.md",
		"user_agreement":  "user_agreement.md",
		"faq":             "faq.md",
	}
	var allChunks []chunk
	for cat, file := range policyFiles {
		data, err := os.ReadFile(filepath.Join("knowledge/policies", file))
		if err != nil { log.Printf("skip %s: %v", file, err); continue }
		chunks := chunkPolicy(string(data), cat)
		allChunks = append(allChunks, chunks...)
		fmt.Printf("  %s → %d chunks\n", cat, len(chunks))
	}
	fmt.Printf("Total chunks: %d\n", len(allChunks))

	// Embed and insert in batches
	apiKey := strings.TrimSpace(os.Getenv("SILICONFLOW_API_KEY"))
	if apiKey == "" { log.Fatal("SILICONFLOW_API_KEY not set") }

	inserted := 0
	for i := 0; i < len(allChunks); i += batchSize {
		end := i + batchSize
		if end > len(allChunks) { end = len(allChunks) }
		batch := allChunks[i:end]

		texts := make([]string, len(batch))
		for j, ch := range batch { texts[j] = ch.content }

		embeddings, err := getEmbeddings(ctx, apiKey, texts)
		if err != nil { log.Printf("embed error batch %d: %v", i, err); continue }

		contentCol := make([]string, len(batch))
		sourceCol := make([]string, len(batch))
		refIDCol := make([]int64, len(batch))
		categoryCol := make([]string, len(batch))
		vectorCol := make([][]float32, len(batch))

		for j, ch := range batch {
			contentCol[j] = truncate(ch.content, 4096)
			sourceCol[j] = ch.source
			refIDCol[j] = 0
			categoryCol[j] = ch.category
			vectorCol[j] = embeddings[j]
		}

		columns := []entity.Column{
			entity.NewColumnVarChar("content", contentCol),
			entity.NewColumnVarChar("source", sourceCol),
			entity.NewColumnInt64("ref_id", refIDCol),
			entity.NewColumnVarChar("category", categoryCol),
			entity.NewColumnFloatVector("vector", embeddingDim, vectorCol),
		}

		if _, err := c.Insert(ctx, "knowledge_chunks", "", columns...); err != nil {
			log.Printf("insert error batch %d: %v", i, err); continue
		}
		inserted += len(batch)
		fmt.Printf("  Inserted %d/%d\n", inserted, len(allChunks))
	}

	c.Flush(ctx, "knowledge_chunks", false)

	// Verify
	stats, _ := c.GetCollectionStatistics(ctx, "knowledge_chunks")
	fmt.Printf("\nDone! Total rows in Milvus: %s\n", stats["row_count"])
}

func getEmbeddings(ctx context.Context, apiKey string, texts []string) ([][]float32, error) {
	reqBody, _ := json.Marshal(embedReq{Input: texts, Model: embeddingModel})
	req, _ := http.NewRequestWithContext(ctx, "POST", siliconFlowURL, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 { return nil, fmt.Errorf("API %d: %s", resp.StatusCode, string(body)) }

	var result embedResp
	json.Unmarshal(body, &result)
	embeddings := make([][]float32, len(texts))
	for _, d := range result.Data { embeddings[d.Index] = d.Embedding }
	return embeddings, nil
}

func chunkPolicy(content, category string) []chunk {
	var chunks []chunk
	paragraphs := strings.Split(content, "\n\n")
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" || utf8.RuneCountInString(para) < minChunkLen { continue }
		if strings.HasPrefix(para, "#") && !strings.Contains(para, "\n") { continue }
		runes := []rune(para)
		if len(runes) <= chunkSize {
			chunks = append(chunks, chunk{content: para, source: "policy", category: category})
		} else {
			step := chunkSize - chunkOverlap
			for i := 0; i < len(runes); {
				end := i + chunkSize
				if end > len(runes) { end = len(runes) }
				ch := strings.TrimSpace(string(runes[i:end]))
				if utf8.RuneCountInString(ch) >= minChunkLen {
					chunks = append(chunks, chunk{content: ch, source: "policy", category: category})
				}
				if end >= len(runes) { break }
				i += step
			}
		}
	}
	return chunks
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n { return string(runes[:n]) }
	return s
}
