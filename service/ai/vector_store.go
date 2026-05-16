package ai

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	milvusCollectionName = "knowledge_chunks"
	milvusVectorField     = "vector"
	milvusContentField    = "content"
	milvusSourceField     = "source"
	milvusRefIDField      = "ref_id"
	milvusCategoryField   = "category"
	dimField              = "dim"
)

var (
	milvusClient client.Client
	milvusMu     sync.Mutex
)

// SetMilvusClient injects the Milvus client. Called once at startup.
func SetMilvusClient(c client.Client) {
	milvusClient = c
}

// InitMilvusCollection creates the knowledge_chunks collection if it doesn't exist.
// Schema: [id:int64 PK] [content:VARCHAR] [source:VARCHAR] [ref_id:INT64]
// [category:VARCHAR] [vector:FLOAT_VECTOR(1024)]
func InitMilvusCollection(ctx context.Context) error {
	milvusMu.Lock()
	defer milvusMu.Unlock()

	if milvusClient == nil {
		return fmt.Errorf("Milvus client not initialized")
	}

	// Check if collection already exists
	has, err := milvusClient.HasCollection(ctx, milvusCollectionName)
	if err != nil {
		return fmt.Errorf("检查 collection 失败: %w", err)
	}
	if has {
		return nil // already exists
	}

	// Build schema
	schema := entity.NewSchema().
		WithName(milvusCollectionName).
		WithDescription("RAG knowledge chunks").
		WithField(entity.NewField().WithName("id").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
		WithField(entity.NewField().WithName(milvusContentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(4096)).
		WithField(entity.NewField().WithName(milvusSourceField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(64)).
		WithField(entity.NewField().WithName(milvusRefIDField).WithDataType(entity.FieldTypeInt64)).
		WithField(entity.NewField().WithName(milvusCategoryField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(128)).
		WithField(entity.NewField().WithName(milvusVectorField).WithDataType(entity.FieldTypeFloatVector).WithDim(embeddingDim))

	if err := milvusClient.CreateCollection(ctx, schema, 2); err != nil {
		return fmt.Errorf("创建 collection 失败: %w", err)
	}

	// Create IVF_FLAT index
	idx, err := entity.NewIndexIvfFlat(entity.L2, 128)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	if err := milvusClient.CreateIndex(ctx, milvusCollectionName, milvusVectorField, idx, false); err != nil {
		return fmt.Errorf("创建向量索引失败: %w", err)
	}

	// Load collection into memory
	if err := milvusClient.LoadCollection(ctx, milvusCollectionName, false); err != nil {
		return fmt.Errorf("加载 collection 失败: %w", err)
	}

	log.Printf("[Milvus] collection %s created with %d-dim vectors", milvusCollectionName, embeddingDim)
	return nil
}

// StoreChunks inserts chunks and their embeddings into Milvus.
func StoreChunks(ctx context.Context, chunks []Chunk, embeddings [][]float32) (int, error) {
	milvusMu.Lock()
	defer milvusMu.Unlock()

	if milvusClient == nil {
		return 0, fmt.Errorf("Milvus client not initialized")
	}

	if len(chunks) != len(embeddings) {
		return 0, fmt.Errorf("chunks 与 embeddings 数量不匹配: %d != %d", len(chunks), len(embeddings))
	}

	if len(chunks) == 0 {
		return 0, nil
	}

	// Build column data (id is auto-generated, skip it)
	contentCol := make([]string, len(chunks))
	sourceCol := make([]string, len(chunks))
	refIDCol := make([]int64, len(chunks))
	categoryCol := make([]string, len(chunks))
	vectorCol := make([][]float32, len(chunks))

	for i := range chunks {
		contentCol[i] = truncateString(chunks[i].Content, 4096)
		sourceCol[i] = chunks[i].Source
		refIDCol[i] = int64(chunks[i].RefID)
		categoryCol[i] = chunks[i].Category
		vectorCol[i] = embeddings[i]
	}

	// Create columns (omit id — auto-generated)
	columns := []entity.Column{
		entity.NewColumnVarChar(milvusContentField, contentCol),
		entity.NewColumnVarChar(milvusSourceField, sourceCol),
		entity.NewColumnInt64(milvusRefIDField, refIDCol),
		entity.NewColumnVarChar(milvusCategoryField, categoryCol),
		entity.NewColumnFloatVector(milvusVectorField, embeddingDim, vectorCol),
	}

	_, err := milvusClient.Insert(ctx, milvusCollectionName, "", columns...)
	if err != nil {
		return 0, fmt.Errorf("插入数据失败: %w", err)
	}

	// Flush to make data searchable
	if err := milvusClient.Flush(ctx, milvusCollectionName, false); err != nil {
		log.Printf("[Milvus] flush warning: %v", err)
	}

	return len(chunks), nil
}

// SearchSimilar performs vector similarity search in Milvus.
func SearchSimilar(ctx context.Context, queryVec []float32, maxResults int, sourceFilter string) ([]ChunkWithScore, error) {
	milvusMu.Lock()
	defer milvusMu.Unlock()

	if milvusClient == nil {
		return nil, fmt.Errorf("Milvus client not initialized")
	}

	if maxResults <= 0 {
		maxResults = 5
	}

	sp, _ := entity.NewIndexIvfFlatSearchParam(16)

	// Build expression filter
	expr := ""
	if sourceFilter != "" {
		expr = fmt.Sprintf(`%s == "%s"`, milvusSourceField, sourceFilter)
	}

	vectors := []entity.Vector{entity.FloatVector(queryVec)}

	results, err := milvusClient.Search(
		ctx,
		milvusCollectionName,
		nil,        // partition names
		expr,       // expression filter
		[]string{milvusContentField, milvusSourceField, milvusRefIDField, milvusCategoryField},
		vectors,
		milvusVectorField,
		entity.L2,  // metric type
		maxResults,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("向量搜索失败: %w", err)
	}

	chunks := make([]ChunkWithScore, 0, len(results))
	for _, sr := range results {
		for i := 0; i < sr.ResultCount; i++ {
			content, _ := sr.Fields.GetColumn(milvusContentField).GetAsString(i)
			source, _ := sr.Fields.GetColumn(milvusSourceField).GetAsString(i)
			category, _ := sr.Fields.GetColumn(milvusCategoryField).GetAsString(i)
			refID, _ := sr.Fields.GetColumn(milvusRefIDField).GetAsInt64(i)

			// L2 distance -> similarity score (closer = smaller L2 = higher score)
			score := 1.0 / (1.0 + float64(sr.Scores[i]))

			chunks = append(chunks, ChunkWithScore{
				Chunk: Chunk{
					Content:  content,
					Source:   source,
					RefID:    uint(refID),
					Category: category,
				},
				Score: score,
			})
		}
	}

	return chunks, nil
}

// GetChunkCount returns the number of entities in the collection.
func GetChunkCount(ctx context.Context) (int64, error) {
	if milvusClient == nil {
		return 0, fmt.Errorf("Milvus client not initialized")
	}

	stats, err := milvusClient.GetCollectionStatistics(ctx, milvusCollectionName)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(stats["row_count"], 10, 64)
}

// ClearAllChunks drops and recreates the collection.
func ClearAllChunks(ctx context.Context) error {
	if milvusClient == nil {
		return fmt.Errorf("Milvus client not initialized")
	}

	has, err := milvusClient.HasCollection(ctx, milvusCollectionName)
	if err != nil {
		return err
	}

	if has {
		if err := milvusClient.DropCollection(ctx, milvusCollectionName); err != nil {
			return fmt.Errorf("删除 collection 失败: %w", err)
		}
	}

	return InitMilvusCollection(ctx)
}

// ClearChunksBySource removes chunks by source expression.
func ClearChunksBySource(ctx context.Context, source string) (int, error) {
	if milvusClient == nil {
		return 0, fmt.Errorf("Milvus client not initialized")
	}

	expr := fmt.Sprintf(`%s == "%s"`, milvusSourceField, source)
	if err := milvusClient.Delete(ctx, milvusCollectionName, "", expr); err != nil {
		return 0, fmt.Errorf("删除失败: %w", err)
	}

	return 0, nil
}

// ListChunksBySource returns chunk previews filtered by source.
func ListChunksBySource(ctx context.Context, source string) ([]Chunk, error) {
	if milvusClient == nil {
		return nil, fmt.Errorf("Milvus client not initialized")
	}

	// Use query to get chunks
	expr := fmt.Sprintf(`%s == "%s"`, milvusSourceField, source)
	if source == "" {
		expr = "id >= 0" // match all
	}

	resultSet, err := milvusClient.Query(
		ctx,
		milvusCollectionName,
		nil,
		expr,
		[]string{milvusContentField, milvusSourceField, milvusRefIDField, milvusCategoryField},
		client.WithLimit(100),
	)
	if err != nil {
		return nil, err
	}

	contentCol := resultSet.GetColumn(milvusContentField)
	sourceCol := resultSet.GetColumn(milvusSourceField)
	categoryCol := resultSet.GetColumn(milvusCategoryField)
	refIDCol := resultSet.GetColumn(milvusRefIDField)

	n := 0
	for _, col := range resultSet {
		if col != nil && col.Len() > n {
			n = col.Len()
		}
	}

	var chunks []Chunk
	for i := 0; i < n; i++ {
		content, _ := contentCol.GetAsString(i)
		source, _ := sourceCol.GetAsString(i)
		category, _ := categoryCol.GetAsString(i)
		refID, _ := refIDCol.GetAsInt64(i)

		// Truncate for preview
		contentRunes := []rune(content)
		if len(contentRunes) > 100 {
			content = string(contentRunes[:100]) + "..."
		}

		chunks = append(chunks, Chunk{
			Content:  content,
			Source:   source[0:intMin(10, len(source))],
			RefID:    uint(refID),
			Category: category[0:intMin(30, len(category))],
		})
	}

	return chunks, nil
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Ensure collection is loaded on startup.
func EnsureCollectionLoaded(ctx context.Context) error {
	if milvusClient == nil {
		return fmt.Errorf("Milvus client not initialized")
	}
	return milvusClient.LoadCollection(ctx, milvusCollectionName, false)
}
