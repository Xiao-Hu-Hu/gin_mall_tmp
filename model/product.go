package model

import (
	"context"
	"gin_mall_tmp/cache"
	"gorm.io/gorm"
	"strconv"
)

type Product struct {
	gorm.Model
	Name          string
	CategoryId    uint
	Title         string
	Info          string
	ImgPath       string `gorm:"not null"`
	Price         string
	DiscountPrice string
	OnSale        bool `gorm:default:false`
	Num           int
	BossID        uint `gorm:"not null"`
	BossName      string
	BossAvatar    string
}

func (product *Product) View() uint64 {
	countStr, _ := cache.RedisClient.Get(context.Background(), cache.ProductViewKey(product.ID)).Result()
	count, _ := strconv.ParseInt(countStr, 10, 64)
	return uint64(count)
}

func (product *Product) AddView() {
	// 增加商品点击数
	cache.RedisClient.Incr(context.Background(), cache.ProductViewKey(product.ID))
	cache.RedisClient.ZIncrBy(context.Background(), cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}
