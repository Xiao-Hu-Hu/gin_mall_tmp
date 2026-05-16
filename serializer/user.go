package serializer

import (
	"gin_mall_tmp/conf"
	"gin_mall_tmp/model"
)

type User struct { // vo : view object
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	Type      int    `json:"type"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:        user.ID,
		UserName:  user.UserName,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Avatar:    conf.Host + conf.HttpPort + conf.AvatarPath + user.Avatar,
		CreatedAt: user.CreatedAt.Unix(),
	}
}
