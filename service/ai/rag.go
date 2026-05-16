package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

const (
	maxRAGResults   = 5
	minScore        = 0.3 // threshold to filter low-similarity results
	maxContextChars = 3000
)

// consultKeywords are used by IsConsultQuery for simple intent detection.
var consultKeywords = []string{
	"退货", "退款", "换货", "配送", "运费", "快递", "物流",
	"支付", "余额", "充值", "发票", "优惠券", "注册", "注销",
	"密码", "账号", "账户", "售后", "投诉", "秒杀", "包邮",
	"客服", "协议", "隐私", "发货", "签收",
}

// RetrieveKnowledge embeds the query and searches Milvus for similar chunks.
func RetrieveKnowledge(ctx context.Context, query string, maxResults int, sourceFilter string) ([]ChunkWithScore, error) {
	if maxResults <= 0 {
		maxResults = maxRAGResults
	}

	embeddings, err := GetEmbeddings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("embedding query failed: %w", err)
	}
	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return nil, fmt.Errorf("empty embedding returned for query")
	}

	results, err := SearchSimilar(ctx, embeddings[0], maxResults, sourceFilter)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// Filter by minimum score
	filtered := make([]ChunkWithScore, 0, len(results))
	for _, r := range results {
		if r.Score >= minScore {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

// BuildRAGPrompt retrieves relevant knowledge and formats it as a prompt context.
// Returns empty string if no relevant chunks found.
func BuildRAGPrompt(ctx context.Context, query string) string {
	chunks, err := RetrieveKnowledge(ctx, query, maxRAGResults, "")
	if err != nil {
		log.Printf("[RAG] retrieval failed: %v", err)
		return ""
	}

	if len(chunks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("以下是商城知识库中与用户问题相关的信息，请严格据此回答：\n\n")

	totalChars := 0
	included := 0
	for _, c := range chunks {
		content := c.Chunk.Content
		charCount := utf8.RuneCountInString(content)
		if totalChars+charCount > maxContextChars {
			break
		}
		sb.WriteString(fmt.Sprintf("【来源：%s】%s\n\n", c.Chunk.Category, content))
		totalChars += charCount
		included++
	}

	if included == 0 {
		return ""
	}

	sb.WriteString("请基于以上知识库内容回答用户问题。如果知识库中没有相关信息，请如实告知用户。")
	return sb.String()
}

// IsConsultQuery checks if the user query is about platform policies or knowledge.
func IsConsultQuery(query string) bool {
	lower := strings.ToLower(query)
	for _, kw := range consultKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// BuildKnowledgeSummary returns a condensed summary of all knowledge topics for the system prompt.
func BuildKnowledgeSummary() string {
	return `商城知识库涵盖以下主题：
- 退货退款：7天无理由退货，15天质量问题退货，换货政策，退款时效
- 配送运费：满99包邮，不满50收12元，50-99收8元，1-7工作日配送
- 支付规则：支持微信/支付宝/银行卡/余额，支付限额，余额充值规则
- 用户协议：注册注销，权利义务，隐私保护，违规处理
- 常见问题：账户、商品、订单、物流、售后、支付、秒杀等 FAQ
- 秒杀规则：特殊价格商品，不支持优惠券和退货，15分钟内支付
- 客服时间：每日9:00-21:00，投诉邮箱 support@mall.com

当用户咨询以上问题时，请调用 search_knowledge 工具检索知识库获取详细内容。`
}
