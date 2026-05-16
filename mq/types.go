package mq

// PayMessage 支付消息
type PayMessage struct {
	OrderID   uint    `json:"order_id"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	BossID    uint    `json:"boss_id"`
	Money     float64 `json:"money"`
	OrderNum  string  `json:"order_num"`
	Num       int     `json:"num"`
}

// InventoryMessage 库存消息
type InventoryMessage struct {
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Num       int  `json:"num"`
}

// OrderCancelMessage 取消订单消息
type OrderCancelMessage struct {
	OrderID uint `json:"order_id"`
	UserID  uint `json:"user_id"`
}
