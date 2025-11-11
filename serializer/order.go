package serializer

import (
	"context"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
)

type Order struct {
	Id            uint    `json:"id" form:"id"`
	OrderNum      string  `json:"order_num" form:"order_num"`
	ProductId     uint    `json:"product_id" form:"product_id"`
	ProductName   string  `json:"product_name" form:"product_name"`
	ImgPath       string  `json:"img_path" form:"img_path"`
	Num           int     `json:"num" form:"num"`
	Price         string  `json:"money" form:"money"`
	DiscountPrice string  `json:"discount_price" form:"discount_price"`
	TotalMoney    float64 `json:"total_money" form:"total_money"`
	AddressId     uint    `json:"address_id" form:"address_id"`
	AddressName   string  `json:"address_name" form:"address_name"`
	AddressPhone  string  `json:"address_phone" form:"address_phone"`
	Address       string  `json:"address" form:"address"`
	UserId        uint    `json:"user_id" form:"user_id"`
	UserNickname  string  `json:"user_nickname" form:"user_nickname"`
	BossId        uint    `json:"boss_id" form:"boss_id"`
	BossNickname  string  `json:"boss_nickname" form:"boss_nickname"`
	Type          uint    `json:"type" form:"type"`
	CreateAt      int64   `json:"create_at" form:"create_at"`
}

func BuildOrder(order *model.Order, product *model.Product, boss *model.User, user *model.User, address *model.Address) Order {
	return Order{
		Id:            order.ID,
		OrderNum:      order.OrderNum,
		ProductId:     order.ProductID,
		ProductName:   product.Name,
		ImgPath:       conf.Host + conf.HttpPort + conf.ProductPath + product.ImgPath,
		Num:           order.Num,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		TotalMoney:    order.Money,
		AddressId:     order.AddressId,
		AddressName:   address.Name,
		AddressPhone:  address.Phone,
		Address:       address.Address,
		UserId:        order.UserID,
		UserNickname:  user.NickName,
		BossId:        order.BossID,
		BossNickname:  boss.NickName,
		Type:          order.Type,
		CreateAt:      order.CreatedAt.Unix(),
	}
}

func BuildOrders(ctx context.Context, items []*model.Order) (orders []Order) {
	productDao := dao.NewProductDao(ctx)
	bossDao := dao.NewUserDao(ctx)
	userDao := dao.NewUserDao(ctx)
	addressDao := dao.NewAddressDao(ctx)
	for _, item := range items {
		product, err := productDao.GetProductById(item.ProductID)
		if err != nil {
			continue
		}
		boss, err := bossDao.GetBossByProductId(item.ProductID)
		if err != nil {
			continue
		}
		user, err := userDao.GetUserById(item.UserID)
		if err != nil {
			continue
		}
		address, err := addressDao.GetAddressByAid(item.AddressId, item.UserID)

		order := BuildOrder(item, product, boss, user, address)
		orders = append(orders, order)
	}
	return orders
}
