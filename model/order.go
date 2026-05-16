package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	BossID    uint `gorm:"not null"`
	AddressId uint `gorm:"not null"`
	Num       int
	OrderNum  string // 订单号
	Type      uint   // 1 未支付 2 已支付
	Money     float64
}
