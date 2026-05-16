package ai

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleUser      ChatRole = "user"
	RoleAssistant ChatRole = "assistant"
)

const (
	IntentChat        = "chat"
	IntentNeed        = "shopping_need"
	IntentConsult     = "consult"
	IntentCreateOrder = "create_order"
)

type ChatMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

type CustomerServiceRequest struct {
	SessionID string        `json:"session_id" form:"session_id"`
	Message   string        `json:"message" form:"message" binding:"required"`
	History   []ChatMessage `json:"history"`
}

type searchProductsInput struct {
	Query string `json:"query" jsonschema:"description=用户需求关键词，例如电脑、手机、拍照好、轻薄本"`
	Limit int    `json:"limit" jsonschema:"description=返回商品条数，建议 1 到 5"`
}

type createOrderInput struct {
	ProductID uint `json:"product_id" jsonschema:"description=商品 ID"`
	Num       int  `json:"num" jsonschema:"description=购买数量，必须大于 0"`
	AddressID uint `json:"address_id" jsonschema:"description=用户收货地址 ID"`
}

// Chunk represents a single knowledge chunk for vector storage.
type Chunk struct {
	Content  string `json:"content"`
	Source   string `json:"source"`
	RefID    uint   `json:"ref_id"`
	Category string `json:"category"`
}

// ChunkWithScore is a search result with similarity score.
type ChunkWithScore struct {
	Chunk Chunk   `json:"chunk"`
	Score float64 `json:"score"`
}

// KnowledgeSyncRequest is the request to sync knowledge (policies + products).
type KnowledgeSyncRequest struct {
	Target string `json:"target" form:"target"` // "policies", "products", or "" for all
}

// KnowledgeSyncPayload is the response for knowledge sync.
type KnowledgeSyncPayload struct {
	PolicyCount  int   `json:"policy_count"`
	ProductCount int   `json:"product_count"`
	TotalChunks  int64 `json:"total_chunks"`
}

// KnowledgeSearchRequest is the request to search knowledge base.
type KnowledgeSearchRequest struct {
	Query  string `json:"query" form:"query" binding:"required"`
	Source string `json:"source" form:"source"` // optional source filter
	Limit  int    `json:"limit" form:"limit"`    // max results, default 5
}
