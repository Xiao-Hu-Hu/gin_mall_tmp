package model

import (
	"gorm.io/gorm"
	"time"
)

// SeckillActivity 秒杀活动
type SeckillActivity struct {
	gorm.Model
	ProductID    uint      `gorm:"index;not null"` // 关联商品ID
	SeckillPrice float64   `gorm:"not null"`        // 秒杀价
	TotalStock   int       `gorm:"not null"`        // 秒杀总库存
	StartAt      time.Time `gorm:"not null"`        // 秒杀开始时间
	EndAt        time.Time `gorm:"not null"`        // 秒杀结束时间
	Status       uint      `gorm:"default:0"`       // 0=未开始, 1=进行中, 2=已结束
}

// SeckillOrder 秒杀订单
type SeckillOrder struct {
	gorm.Model
	UserID     uint    `gorm:"index;not null"`  // 用户ID
	ActivityID uint    `gorm:"index;not null"`  // 秒杀活动ID
	ProductID  uint    `gorm:"not null"`         // 商品ID
	OrderNum   string  `gorm:"uniqueIndex;not null"` // 订单号
	Money      float64 `gorm:"not null"`         // 秒杀价
	Status     uint    `gorm:"default:0"`        // 0=待支付, 1=已支付, 2=已取消
}
