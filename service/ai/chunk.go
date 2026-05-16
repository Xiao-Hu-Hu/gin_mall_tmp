package ai

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"gin_mall_tmp/model"
)

const (
	chunkSize    = 500 // max characters per chunk
	chunkOverlap = 50  // overlap between consecutive chunks
	minChunkLen  = 10  // minimum characters for a valid chunk
)

// ChunkText splits a long text into overlapping chunks of approximately chunkSize characters.
func ChunkText(text string) []string {
	runes := []rune(text)
	total := len(runes)

	if total <= chunkSize {
		return []string{text}
	}

	var chunks []string
	step := chunkSize - chunkOverlap
	for i := 0; i < total; {
		end := i + chunkSize
		if end > total {
			end = total
		}

		chunk := strings.TrimSpace(string(runes[i:end]))
		if utf8.RuneCountInString(chunk) >= minChunkLen {
			chunks = append(chunks, chunk)
		}

		if end >= total {
			break
		}
		i += step
	}

	return chunks
}

// ChunkProduct converts a product record into knowledge chunks for vector search.
func ChunkProduct(product *model.Product) []Chunk {
	text := fmt.Sprintf(
		"商品名称：%s\n标题：%s\n描述：%s\n价格：%s元\n折扣价：%s元\n库存：%d件\n分类ID：%d",
		product.Name, product.Title, product.Info,
		product.Price, product.DiscountPrice, product.Num, product.CategoryId,
	)

	pieces := ChunkText(text)
	chunks := make([]Chunk, 0, len(pieces))
	for _, p := range pieces {
		chunks = append(chunks, Chunk{
			Content:  p,
			Source:   "product",
			RefID:    product.ID,
			Category: fmt.Sprintf("category_%d", product.CategoryId),
		})
	}
	return chunks
}

// ChunkPolicy splits a Markdown policy document into chunks.
func ChunkPolicy(content string, policyName string) []Chunk {
	var chunks []Chunk

	paragraphs := strings.Split(content, "\n\n")

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" || utf8.RuneCountInString(para) < minChunkLen {
			continue
		}

		if strings.HasPrefix(para, "#") && !strings.Contains(para, "\n") {
			continue
		}

		if utf8.RuneCountInString(para) > chunkSize {
			subChunks := ChunkText(para)
			for _, sc := range subChunks {
				if utf8.RuneCountInString(sc) >= minChunkLen {
					chunks = append(chunks, Chunk{
						Content:  sc,
						Source:   "policy",
						Category: policyName,
					})
				}
			}
		} else {
			chunks = append(chunks, Chunk{
				Content:  para,
				Source:   "policy",
				Category: policyName,
			})
		}
	}

	return chunks
}
