package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	User      User    `gorm:"foreignkey:UserId"`
	UserId    uint    `gorm:"not null"`
	Product   Product `gorm:"foreignkey:ProductId"`
	ProductId uint    `gorm:"not null"`
	Boss      User    `gorm:"foreignkey:BossId"`
	BossId    uint    `gorm:"not null"`
}
