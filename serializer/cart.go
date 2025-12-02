package serializer

import (
	"context"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
)

type Cart struct {
	Id            uint   `json:"id"`
	UserId        uint   `json:"user_id"`
	ProductId     uint   `json:"product_id"`
	CreatAt       int64  `json:"creatAt"`
	Name          string `json:"name"`
	Num           int    `json:"num"`
	MaxNum        int    `json:"maxNum"`
	ImgPath       string `json:"img_path"`
	Check         bool   `json:"check"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
}

func BuildCart(cart *model.Cart, product *model.Product, boss *model.User) Cart {
	return Cart{
		Id:            cart.ID,
		UserId:        cart.UserId,
		ProductId:     product.ID,
		CreatAt:       cart.CreatedAt.Unix(),
		Num:           int(cart.Num),
		MaxNum:        int(cart.MaxNum),
		Check:         cart.Check,
		Name:          product.Name,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		BossId:        boss.ID,
		BossName:      boss.UserName,
	}
}

func BuildCarts(ctx context.Context, items []*model.Cart) (carts []Cart) {
	productDao := dao.NewProductDao(ctx)
	bossDao := dao.NewUserDao(ctx)
	for _, item := range items {
		// 跳过无效的商品ID（0或空值）
		if item.ProductId == 0 {
			continue
		}
		product, err := productDao.GetProductById(item.ProductId)
		if err != nil {
			continue
		}
		// 跳过无效的老板ID
		if item.BossId == 0 {
			continue
		}
		boss, err := bossDao.GetUserById(item.BossId)
		if err != nil {
			continue
		}
		cart := BuildCart(item, product, boss)
		carts = append(carts, cart)
	}
	return carts
}
