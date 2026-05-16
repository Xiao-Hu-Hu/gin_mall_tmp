package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
)

const embeddingBatchSize = 20 // max texts per embedding API call

// policyFiles maps display name -> file path relative to knowledge/policies/.
var policyFiles = map[string]string{
	"return_policy":   "return_policy.md",
	"shipping_policy": "shipping_policy.md",
	"payment_policy":  "payment_policy.md",
	"user_agreement":  "user_agreement.md",
	"faq":             "faq.md",
}

// SyncPolicies reads all policy markdown files, chunks, embeds, and stores them in Milvus.
func SyncPolicies(ctx context.Context) (int, error) {
	baseDir := "knowledge/policies"
	var allChunks []Chunk

	for category, filename := range policyFiles {
		path := filepath.Join(baseDir, filename)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[Sync] skip %s: %v", path, err)
			continue
		}

		chunks := ChunkPolicy(string(data), category)
		allChunks = append(allChunks, chunks...)
		log.Printf("[Sync] %s → %d chunks", category, len(chunks))
	}

	if len(allChunks) == 0 {
		return 0, fmt.Errorf("no policy chunks generated")
	}

	// Clear old policy data
	if _, err := ClearChunksBySource(ctx, "policy"); err != nil {
		log.Printf("[Sync] warning: clear old policies failed: %v", err)
	}

	return storeChunksInBatches(ctx, allChunks)
}

// SyncAllProducts chunks all on-sale products, embeds, and stores them in Milvus.
func SyncAllProducts(ctx context.Context) (int, error) {
	db := dao.NewProductDao(ctx)
	var products []*model.Product
	if err := db.DB.Where("on_sale = ?", true).Find(&products).Error; err != nil {
		return 0, fmt.Errorf("query products failed: %w", err)
	}

	if len(products) == 0 {
		return 0, fmt.Errorf("no products found")
	}

	var allChunks []Chunk
	for _, p := range products {
		chunks := ChunkProduct(p)
		allChunks = append(allChunks, chunks...)
	}
	log.Printf("[Sync] %d products → %d chunks", len(products), len(allChunks))

	// Clear old product data
	if _, err := ClearChunksBySource(ctx, "product"); err != nil {
		log.Printf("[Sync] warning: clear old products failed: %v", err)
	}

	return storeChunksInBatches(ctx, allChunks)
}

// SyncKnowledgeBase syncs all knowledge (policies + products) into Milvus.
func SyncKnowledgeBase(ctx context.Context, target string) (policyCount, productCount int, totalChunks int64, err error) {
	if target == "" || target == "policies" {
		policyCount, err = SyncPolicies(ctx)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("sync policies failed: %w", err)
		}
	}

	if target == "" || target == "products" {
		productCount, err = SyncAllProducts(ctx)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("sync products failed: %w", err)
		}
	}

	totalChunks, err = GetChunkCount(ctx)
	if err != nil {
		log.Printf("[Sync] get chunk count failed: %v", err)
	}

	return policyCount, productCount, totalChunks, nil
}

// storeChunksInBatches sends chunks to the embedding API in batches and inserts into Milvus.
func storeChunksInBatches(ctx context.Context, chunks []Chunk) (int, error) {
	total := 0

	for i := 0; i < len(chunks); {
		end := i + embeddingBatchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		batch := chunks[i:end]
		texts := make([]string, len(batch))
		for j, c := range batch {
			texts[j] = c.Content
		}

		embeddings, err := GetEmbeddings(ctx, texts)
		if err != nil {
			return total, fmt.Errorf("embedding batch [%d:%d] failed: %w", i, end, err)
		}

		n, err := StoreChunks(ctx, batch, embeddings)
		if err != nil {
			return total, fmt.Errorf("store batch [%d:%d] failed: %w", i, end, err)
		}
		total += n

		log.Printf("[Sync] stored batch %d/%d (%d chunks)", end, len(chunks), n)
		i = end
	}

	return total, nil
}

// readPolicyFile reads a single policy file and returns its content.
func readPolicyFile(category string) string {
	baseDir := "knowledge/policies"
	filename, ok := policyFiles[category]
	if !ok {
		filename = category + ".md"
	}

	data, err := os.ReadFile(filepath.Join(baseDir, filename))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
