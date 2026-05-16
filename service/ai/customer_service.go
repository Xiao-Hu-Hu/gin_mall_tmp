package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin_mall_tmp/cache"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/serializer"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	openai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

const (
	defaultDeepSeekBaseURL = "https://api.deepseek.com"
	defaultDeepSeekModel   = "deepseek-chat"
	orderAPIPath           = "/api/v1/orders"
)

type CustomerService struct {
	request   *CustomerServiceRequest
	userID    uint
	userToken string
}

type toolSessionState struct {
	intent           string
	recommendations  []serializer.AIRecommendedProduct
	addresses        []serializer.Address
	orderResult      *serializer.AIOrderResult
	knowledgeChunks  []ChunkWithScore
}

type searchKnowledgeInput struct {
	Query string `json:"query" jsonschema:"description=用户在知识库中搜索的关键词，例如退货政策、运费标准"`
}

type orderAPIResponse struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func NewCustomerService(request *CustomerServiceRequest, userID uint, userToken string) *CustomerService {
	return &CustomerService{
		request:   request,
		userID:    userID,
		userToken: userToken,
	}
}

func (s *CustomerService) Chat(ctx context.Context) serializer.Response {
	if strings.TrimSpace(s.request.Message) == "" {
		return serializer.Response{Status: http.StatusBadRequest, Msg: "参数错误", Error: "message 不能为空"}
	}

	// 获取或创建会话 ID
	sessionID := s.request.SessionID
	if sessionID == "" {
		sessionID = generateSessionID(s.userID)
	}

	// 从 Redis 加载历史上下文
	sessionCtx, err := cache.GetAISession(sessionID)
	if err != nil {
		return serializer.Response{Status: http.StatusInternalServerError, Msg: "加载会话失败", Error: err.Error()}
	}

	modelClient, err := s.newChatModel(ctx)
	if err != nil {
		return serializer.Response{Status: http.StatusInternalServerError, Msg: "AI 客服初始化失败", Error: err.Error()}
	}

	state := &toolSessionState{intent: IntentChat}
	tools, err := s.buildTools(state)
	if err != nil {
		return serializer.Response{Status: http.StatusInternalServerError, Msg: "AI 工具初始化失败", Error: err.Error()}
	}

	agentIns, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: modelClient,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
		MaxStep: 8,
	})
	if err != nil {
		return serializer.Response{Status: http.StatusInternalServerError, Msg: "AI Agent 初始化失败", Error: err.Error()}
	}

	// 构建消息（使用 Redis 中的历史 + 当前消息）
	messages := s.buildMessagesWithContext(sessionCtx)

	reply, err := agentIns.Generate(ctx, messages)
	if err != nil {
		return serializer.Response{Status: http.StatusInternalServerError, Msg: "AI 客服调用失败", Error: err.Error()}
	}

	// 保存用户消息到 Redis
	cache.AppendAIMessage(sessionID, s.userID, cache.AIMessage{
		Role:      "user",
		Content:   s.request.Message,
		Timestamp: time.Now(),
	})

	// 保存 AI 回复到 Redis
	cache.AppendAIMessage(sessionID, s.userID, cache.AIMessage{
		Role:      "assistant",
		Content:   reply.Content,
		Timestamp: time.Now(),
	})

	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "ok",
		Data: serializer.AIChatPayload{
			SessionID:       sessionID,
			Intent:          state.intent,
			Reply:           strings.TrimSpace(reply.Content),
			Recommendations: state.recommendations,
			Addresses:       state.addresses,
			OrderResult:     state.orderResult,
		},
	}
}

func (s *CustomerService) newChatModel(ctx context.Context) (*openai.ChatModel, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("未读取到环境变量 DEEPSEEK_API_KEY")
	}

	baseURL := strings.TrimSpace(os.Getenv("DEEPSEEK_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultDeepSeekBaseURL
	}

	modelName := strings.TrimSpace(os.Getenv("DEEPSEEK_MODEL"))
	if modelName == "" {
		modelName = defaultDeepSeekModel
	}

	temperature := float32(0.2)
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       modelName,
		Temperature: &temperature,
		Timeout:     60 * time.Second,
	})
}

func (s *CustomerService) buildMessages() []*schema.Message {
	messages := []*schema.Message{{
		Role: schema.System,
		Content: "你是商城 AI 智能客服。先判断用户是在闲聊、商品咨询还是要求代下单。" +
			"如果只是闲聊，直接回答，不要调用工具。" +
			"如果用户在表达购物需求或让你推荐商品，必须调用 search_products。" +
			"如果用户明确让你帮忙下单，且已提供 product_id、num、address_id，就调用 create_order；若缺字段，先调用 list_user_addresses 或追问。" +
			"禁止编造商品、地址、订单结果。",
	}}

	for _, item := range s.request.History {
		if strings.TrimSpace(item.Content) == "" {
			continue
		}

		role := schema.User
		switch item.Role {
		case RoleSystem:
			role = schema.System
		case RoleAssistant:
			role = schema.Assistant
		case RoleUser:
			role = schema.User
		}

		messages = append(messages, &schema.Message{
			Role:    role,
			Content: item.Content,
		})
	}

	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: s.request.Message,
	})
	return messages
}

func (s *CustomerService) buildMessagesWithContext(sessionCtx *cache.AISessionContext) []*schema.Message {
	messages := []*schema.Message{{
		Role: schema.System,
		Content: "你是商城 AI 智能客服。分析用户意图后选择对应策略：\n" +
			"1. 闲聊/问候 → 直接回答，不要调用工具。\n" +
			"2. 退货/换货/配送/运费/支付/退款/发票/优惠券/秒杀规则/客服时间/账户/注册/注销/投诉等政策咨询 → 必须调用 search_knowledge 后再回答，严格基于知识库内容，不得编造。\n" +
			"3. 购物需求/商品推荐 → 调用 search_products。\n" +
			"4. 明确代下单（已确认 product_id, num, address_id）→ 调用 create_order；缺字段则先调用 list_user_addresses 或追问。\n" +
			"禁止编造商品、地址、订单结果。禁止编造政策规则。",
	}}

	// 从 Redis 上下文加载历史
	if sessionCtx != nil {
		for _, msg := range sessionCtx.History {
			role := schema.User
			if msg.Role == "assistant" {
				role = schema.Assistant
			}
			messages = append(messages, &schema.Message{
				Role:    role,
				Content: msg.Content,
			})
		}
	}

	// 添加当前消息
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: s.request.Message,
	})

	return messages
}

func generateSessionID(userID uint) string {
	return fmt.Sprintf("%d_%d", userID, time.Now().UnixNano())
}

func (s *CustomerService) buildTools(state *toolSessionState) ([]tool.BaseTool, error) {
	searchProductTool, err := einotool.InferTool("search_products", "根据用户购物需求搜索并推荐商城商品", func(ctx context.Context, input searchProductsInput) (string, error) {
		return s.searchProducts(ctx, input, state)
	})
	if err != nil {
		return nil, err
	}

	addressTool, err := einotool.InferTool("list_user_addresses", "查询当前登录用户的收货地址列表，供下单前确认地址 ID", func(ctx context.Context, _ struct{}) (string, error) {
		return s.listUserAddresses(ctx, state)
	})
	if err != nil {
		return nil, err
	}

	orderTool, err := einotool.InferTool("create_order", "当用户明确要求客服代下单时，调用商城现有创建订单接口创建订单", func(ctx context.Context, input createOrderInput) (string, error) {
		return s.createOrderByAPI(ctx, input, state)
	})
	if err != nil {
		return nil, err
	}

	knowledgeTool, err := einotool.InferTool("search_knowledge", "当用户咨询退货/换货/配送/运费/支付/退款/发票/优惠券/秒杀规则/客服时间/账户/投诉等政策问题时，检索知识库获取准确答案", func(ctx context.Context, input searchKnowledgeInput) (string, error) {
		return s.searchKnowledge(ctx, input, state)
	})
	if err != nil {
		return nil, err
	}

	return []tool.BaseTool{searchProductTool, addressTool, orderTool, knowledgeTool}, nil
}

func (s *CustomerService) searchProducts(ctx context.Context, input searchProductsInput, state *toolSessionState) (string, error) {
	query := strings.TrimSpace(input.Query)
	if query == "" {
		return "", errors.New("query 不能为空")
	}

	productDAO := dao.NewProductDao(ctx)
	products, err := productDAO.SearchProductsForAI(query, input.Limit)
	if err != nil {
		return "", err
	}

	recommendations := make([]serializer.AIRecommendedProduct, 0, len(products))
	for _, item := range products {
		recommendations = append(recommendations, buildAIRecommendedProduct(item))
	}

	state.intent = IntentNeed
	state.recommendations = recommendations

	if len(recommendations) == 0 {
		return `{"message":"没有找到匹配商品"}`, nil
	}

	body, err := json.Marshal(recommendations)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *CustomerService) listUserAddresses(ctx context.Context, state *toolSessionState) (string, error) {
	addressDAO := dao.NewAddressDao(ctx)
	addresses, err := addressDAO.ListAddressByUserId(s.userID)
	if err != nil {
		return "", err
	}

	result := serializer.BuildAddresses(addresses)
	state.addresses = result

	body, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *CustomerService) createOrderByAPI(ctx context.Context, input createOrderInput, state *toolSessionState) (string, error) {
	if input.ProductID == 0 {
		return "", errors.New("product_id 不能为空")
	}
	if input.Num <= 0 {
		return "", errors.New("num 必须大于 0")
	}
	if input.AddressID == 0 {
		return "", errors.New("address_id 不能为空")
	}

	state.intent = IntentCreateOrder

	payload, err := json.Marshal(map[string]interface{}{
		"product_id": input.ProductID,
		"num":        input.Num,
		"address_id": input.AddressID,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.orderAPIURL(), bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", s.userToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result orderAPIResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析订单接口响应失败: %w", err)
	}

	state.orderResult = &serializer.AIOrderResult{
		Status: result.Status,
		Msg:    result.Msg,
		Error:  result.Error,
		Data:   result.Data,
	}

	return string(body), nil
}

func (s *CustomerService) orderAPIURL() string {
	baseURL := strings.TrimSpace(os.Getenv("AI_ORDER_API_BASE_URL"))
	if baseURL != "" {
		return strings.TrimRight(baseURL, "/") + orderAPIPath
	}
	return strings.TrimRight(conf.Host+conf.HttpPort, "/") + orderAPIPath
}

func (s *CustomerService) searchKnowledge(ctx context.Context, input searchKnowledgeInput, state *toolSessionState) (string, error) {
	query := strings.TrimSpace(input.Query)
	if query == "" {
		return "", errors.New("query 不能为空")
	}

	chunks, err := RetrieveKnowledge(ctx, query, maxRAGResults, "")
	if err != nil {
		return "", fmt.Errorf("知识库检索失败: %w", err)
	}

	state.intent = IntentConsult
	state.knowledgeChunks = chunks

	if len(chunks) == 0 {
		return "知识库中未找到相关信息，请如实告知用户暂无相关资料。", nil
	}

	// Format results
	var sb strings.Builder
	sb.WriteString("知识库检索结果：\n\n")
	for i, c := range chunks {
		sb.WriteString(fmt.Sprintf("--- 相关内容 %d（相似度 %.0f%%）---\n", i+1, c.Score*100))
		sb.WriteString(fmt.Sprintf("来源：%s\n", c.Chunk.Category))
		sb.WriteString(fmt.Sprintf("内容：%s\n\n", c.Chunk.Content))
	}
	sb.WriteString("请严格基于以上知识库内容回答用户，禁止编造。")

	return sb.String(), nil
}

// buildAIRecommendedProduct converts the existing product serializer output
// into a smaller payload dedicated to the AI customer service response.
func buildAIRecommendedProduct(item *model.Product) serializer.AIRecommendedProduct {
	product := serializer.BuildProduct(item)
	return serializer.AIRecommendedProduct{
		Id:            product.Id,
		Name:          product.Name,
		Title:         product.Title,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		ImgPath:       product.ImgPath,
		Stock:         product.Num,
	}
}
