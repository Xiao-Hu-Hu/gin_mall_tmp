package serializer

type AIRecommendedProduct struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	ImgPath       string `json:"img_path"`
	Stock         int    `json:"stock"`
}

type AIOrderResult struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

type AIChatPayload struct {
	SessionID       string                 `json:"session_id"`
	Intent          string                 `json:"intent"`
	Reply           string                 `json:"reply"`
	Recommendations []AIRecommendedProduct `json:"recommendations"`
	Addresses       []Address              `json:"addresses"`
	OrderResult     *AIOrderResult         `json:"order_result,omitempty"`
}
