package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"type:varchar(20);not null"`
	Avatar   string `gorm:"type:varchar(50);not null"`
}
