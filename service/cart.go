package service

import (
	"context"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"strconv"
)

type CartService struct {
	Id        uint `json:"id" form:"id"`
	BossId    uint `json:"boss_id" form:"boss_id"`
	ProductId uint `json:"product_id" form:"product_id"`
	Num       int  `json:"num" form:"num"`
}

func (service *CartService) Create(ctx context.Context, uId uint) serializer.Response {
	var cart *model.Cart
	code := e.Success

	//先判断有没有这个商品
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CartCreate product err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//再判断有没有这个老板
	userDao := dao.NewUserDao(ctx)
	boss, err := userDao.GetUserById(service.BossId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CartCreate boss err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//如果都没问题，那就创建这个订单
	cartDao := dao.NewCartDao(ctx)
	cart = &model.Cart{
		UserId:    uId,
		ProductId: service.ProductId,
		BossId:    service.BossId,
		Num:       uint(service.Num),
	}
	err = cartDao.CreateCart(cart)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreateCart err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildCart(cart, product, boss),
	}
}

func (service *CartService) List(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	cartDao := dao.NewCartDao(ctx)
	cartList, err := cartDao.ListCartByUserId(uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ListCart err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCarts(ctx, cartList), uint(len(cartList)))
}

func (service *CartService) Show(ctx context.Context, uId uint, cId string) serializer.Response {
	code := e.Success
	cartDao := dao.NewCartDao(ctx)
	cartId, _ := strconv.Atoi(cId)
	cart, err := cartDao.ShowCartByCid(uint(cartId), uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ListCart err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	userDao := dao.NewUserDao(ctx)
	boss, _ := userDao.GetUserById(service.BossId)

	productDao := dao.NewProductDao(ctx)
	product, _ := productDao.GetProductById(service.ProductId)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildCart(cart, product, boss),
	}
}

func (service *CartService) Update(ctx context.Context, uId uint, cId string) serializer.Response {
	code := e.Success
	cartDao := dao.NewCartDao(ctx)
	cartId, _ := strconv.Atoi(cId)
	err := cartDao.UpdateCartNumByUserId(uint(cartId), uId, service.Num)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("UpdateCart err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *CartService) Delete(ctx context.Context, uId uint, cId string) serializer.Response {
	cartId, _ := strconv.Atoi(cId)
	code := e.Success
	cartDao := dao.NewCartDao(ctx)
	err := cartDao.DeleteCartByCartId(uint(cartId), uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("DeleteCart err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
