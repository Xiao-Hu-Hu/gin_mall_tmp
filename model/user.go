package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique;not null"`
	Email          string `gorm:"unique;not null"`
	PasswordDigest string `gorm:"not null"`
	NickName       string `gorm:"not null"`
	Status         string
	Avatar         string `gorm:"not null"`
	Money          string
}

const (
	PasswordCost        = 12       //密码加密难度
	Avtive       string = "active" // 激活用户
)

// 密码加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 判断密码是否正确
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
