package serializer

import (
	"context"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
)

type Favorite struct {
	ID            uint   `json:"id"`
	UserId        uint   `json:"user_id"`
	ImgPath       string `json:"img_path"`
	ProductId     uint   `json:"product_id"`
	CreatedAt     int64  `json:"created_at"`
	Name          string `json:"name"`
	CategoryId    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	Num           int    `json:"num"`
	OnSale        bool   `json:"on_sale"`
}

func BuildFavorite(favorite *model.Favorite, product *model.Product, boss *model.User) Favorite {
	return Favorite{
		ID:            favorite.ID,
		UserId:        favorite.UserId,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		ProductId:     favorite.ProductId,
		CreatedAt:     favorite.CreatedAt.Unix(),
		Name:          product.Name,
		CategoryId:    product.CategoryId,
		Title:         product.Title,
		Info:          product.Info,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		BossId:        boss.ID,
		Num:           product.Num,
		OnSale:        product.OnSale,
	}
}

func BuildFavorites(ctx context.Context, items []*model.Favorite) (favorites []Favorite) {
	productDao := dao.NewProductDao(ctx)
	bossDao := dao.NewUserDao(ctx)
	for _, item := range items {
		product, err := productDao.GetProductById(item.ProductId)
		if err != nil {
			continue
		}
		boss, err := bossDao.GetUserById(item.BossId)
		if err != nil {
			continue
		}
		favorite := BuildFavorite(item, product, boss)
		favorites = append(favorites, favorite)
	}
	return favorites
}
