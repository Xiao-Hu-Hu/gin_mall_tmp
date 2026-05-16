// +build ignore

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func main() {
	apiKey := os.Getenv("SILICONFLOW_API_KEY")
	ctx := context.Background()

	c, _ := client.NewClient(ctx, client.Config{Address: "localhost:19530"})
	defer c.Close()

	queries := []string{
		"我都可以用什么方式来支付啊",
		"我该怎么开发票呢",
		"介绍一下平台的退款规则吧",
		"收到质量问题商品能退货吗",
		"运费多少钱",
	}

	for _, q := range queries {
		vecs := embed(ctx, apiKey, []string{q})
		sp, _ := entity.NewIndexIvfFlatSearchParam(16)
		results, _ := c.Search(ctx, "knowledge_chunks", nil, "",
			[]string{"content", "category"},
			[]entity.Vector{entity.FloatVector(vecs[0])},
			"vector", entity.L2, 1, sp)
		if len(results) > 0 && results[0].ResultCount > 0 {
			cat, _ := results[0].Fields.GetColumn("category").GetAsString(0)
			content, _ := results[0].Fields.GetColumn("content").GetAsString(0)
			score := 1.0 / (1.0 + float64(results[0].Scores[0]))
			runes := []rune(content)
			if len(runes) > 80 { content = string(runes[:80]) + "..." }
			fmt.Printf("Q: %-20s → [%.3f] %s: %s\n", q, score, cat, content)
		} else {
			fmt.Printf("Q: %-20s → No results\n", q)
		}
	}
	c.DropCollection(ctx, "knowledge_chunks")
}

func embed(ctx context.Context, apiKey string, texts []string) [][]float32 {
	body, _ := json.Marshal(map[string]interface{}{"input": texts, "model": "Qwen/Qwen3-Embedding-0.6B"})
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.siliconflow.cn/v1/embeddings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, _ := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	b, _ := io.ReadAll(resp.Body)
	var r struct{ Data []struct{ Embedding []float32; Index int } }
	json.Unmarshal(b, &r)
	embs := make([][]float32, len(texts))
	for _, d := range r.Data { embs[d.Index] = d.Embedding }
	return embs
}
