package service

import (
	"context"
	"errors"
	"fmt"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"strconv"
)

type OrderPayService struct {
	OrderId   uint    `json:"order_id" form:"order_id"`
	Money     float64 `json:"money" form:"money"`
	OrderNo   string  `json:"order_no" form:"order_no"`
	ProductId uint    `json:"product_id" form:"product_id"`
	PayTime   string  `json:"pay_time" form:"pay_time"`
	Sign      string  `json:"sign" form:"sign"`
	BossId    uint    `json:"boss_id" form:"boss_id"`
	BossName  string  `json:"boss_name" form:"boss_name"`
	Num       int     `json:"num" form:"num"`
	Key       string  `json:"key" form:"key"` // 支付金额
}

func (service *OrderPayService) PayDown(ctx context.Context, uId uint) serializer.Response {
	util.Encrypt.SetKey(service.Key)
	code := e.Success
	orderDao := dao.NewOrderDao(ctx)
	tx := orderDao.Begin()
	// 根据订单ID获取订单
	order, err := orderDao.ShowOrderByOid(service.OrderId, uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderPay err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	money := order.Money
	userDao := dao.NewUserDao(ctx)
	// 根据UID获取用户，进行金额修改
	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.GetUser err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 用户扣钱
	// 对钱进行解密，减去金额，再加密保存
	moneyStr := util.Encrypt.AesDecoding(user.Money)
	moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)

	if moneyFloat-money < 0.0 {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.GetUser err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  errors.New("金额不足").Error(),
		}
	}
	finMoney := fmt.Sprintf("%f", moneyFloat-money)
	user.Money = util.Encrypt.AesEncoding(finMoney)

	userDao = dao.NewUserDaoByDB(userDao.DB)
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.UpdateUser err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	var boss *model.User
	boss, err = userDao.GetBossByProductId(service.ProductId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.GetBoss err ", err)
		tx.Rollback()
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 商家价钱
	// 解密商家的余额，在增加之后再进行加密保存
	moneyStr = util.Encrypt.AesDecoding(boss.Money)
	moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
	finMoney = fmt.Sprintf("%f", moneyFloat+money)
	boss.Money = util.Encrypt.AesEncoding(finMoney)

	err = userDao.UpdateUserById(boss.ID, boss)
	if err != nil {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.Update err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 对应的商品数目减少
	var product *model.Product
	productDao := dao.NewProductDao(ctx)
	product, err = productDao.GetProductById(service.ProductId)
	if err != nil {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.GetProduct err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if product.Num < service.Num {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("商品数量不足")
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  errors.New("商品数量不足").Error(),
		}
	}
	product.Num = product.Num - order.Num
	err = productDao.UpdateProductById(service.ProductId, product)
	if err != nil {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("UpdateProduct err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 订单状态修改为已支付
	order.Type = 2
	err = orderDao.UpdateOrderByOrderId(service.OrderId, order, uId)
	if err != nil {
		tx.Rollback()
		code = e.Error
		util.LogrusObj.Infoln("OrderPay.Update err ", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 胡传政认为没必要写这个功能，毫无作用
	// 自己的商品+1  同一件商品该怎么办
	//productUser := model.Product{
	//	Name:          product.Name,
	//	CategoryId:    product.CategoryId,
	//	Title:         product.Title,
	//	Info:          product.Info,
	//	ImgPath:       product.ImgPath,
	//	Price:         product.Price,
	//	DiscountPrice: product.DiscountPrice,
	//	OnSale:        false,
	//	Num:           order.Num,
	//	BossID:        uId,
	//	BossName:      user.UserName,
	//	BossAvatar:    user.Avatar,
	//}
	//
	//err = productDao.CreateProduct(&productUser)
	//if err != nil {
	//	tx.Rollback()
	//	code = e.Error
	//	util.LogrusObj.Infoln("OrderPay.CreateProduct err ", err)
	//	return serializer.Response{
	//		Status: code,
	//		Msg:    e.GetMsg(code),
	//		Error:  err.Error(),
	//	}
	//}

	tx.Commit()
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Error:  "",
	}
}
