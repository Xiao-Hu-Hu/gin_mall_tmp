package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// siliconFlowEmbeddingURL is the endpoint for SiliconFlow's Embedding API.
	siliconFlowEmbeddingURL = "https://api.siliconflow.cn/v1/embeddings"

	// defaultEmbeddingModel is the Qwen3-Embedding-0.6B model, 1024-dim.
	defaultEmbeddingModel = "Qwen/Qwen3-Embedding-0.6B"

	// embeddingDim is the output dimension of our chosen embedding model.
	embeddingDim = 1024
)

// embeddingRequest mirrors the OpenAI-compatible request body.
type embeddingRequest struct {
	Input          interface{} `json:"input"` // string or []string
	Model          string      `json:"model"`
	EncodingFormat string      `json:"encoding_format,omitempty"`
}

// embeddingResponse mirrors the OpenAI-compatible response body.
type embeddingResponse struct {
	Object string          `json:"object"`
	Data   []embeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type embeddingData struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// GetEmbeddings converts texts to vector representations via SiliconFlow Embedding API.
//
// Uses the Qwen/Qwen3-Embedding-0.6B model (1024-dim). Supports batch input —
// all texts are sent in a single API call for efficiency.
func GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	apiKey := strings.TrimSpace(os.Getenv("SILICONFLOW_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("未配置 SILICONFLOW_API_KEY 环境变量")
	}

	model := strings.TrimSpace(os.Getenv("EMBEDDING_MODEL"))
	if model == "" {
		model = defaultEmbeddingModel
	}

	reqBody, err := json.Marshal(embeddingRequest{
		Input: texts, // batch: send all texts in one request
		Model: model,
	})
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, siliconFlowEmbeddingURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Embedding API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API 返回错误 (status=%d): %s", resp.StatusCode, string(body))
	}

	var result embeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	embeddings := make([][]float32, len(texts))
	for _, d := range result.Data {
		if d.Index < len(texts) {
			embeddings[d.Index] = d.Embedding
		}
	}

	return embeddings, nil
}
