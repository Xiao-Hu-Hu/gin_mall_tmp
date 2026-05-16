package service

import (
	"context"
	"fmt"
	"gin_mall_tmp/cache"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"math/rand"
	"strconv"
	"time"
)

type OrderService struct {
	ProductId uint `json:"product_id" form:"product_id"`
	Num       int  `json:"num" form:"num"`
	AddressId uint `json:"address_id" form:"address_id"`
	UserId    uint `json:"user_id" form:"user_id"`
	BossId    uint `json:"boss_id" form:"boss_id"`
	OrderNum  uint `json:"order_num" form:"order_num"`
	Type      uint `json:"type" form:"type"`
	model.BasePage
}

func (service *OrderService) Create(ctx context.Context, uId uint) serializer.Response {
	var order *model.Order
	code := e.Success

	//先判断有没有这个商品
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderCreate product err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//再判断有没有这个老板
	userDao := dao.NewUserDao(ctx)
	boss, err := userDao.GetBossByProductId(service.ProductId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderCreate boss err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderCreate user err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 检验用户的地址是否存在
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressByAid(service.AddressId, uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderCreate address err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//获取商品的价格
	money, err := strconv.Atoi(product.Price)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("OrderCreate getProductMoney err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 创建订单号 随机生成的number + 唯一的 product_id + 用户的 ID
	orderNum := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000000))
	productNum := strconv.Itoa(product.Num)
	userNum := strconv.Itoa(int(uId))
	orderNum = orderNum + productNum + userNum

	//如果都没问题，那就创建这个订单
	orderDao := dao.NewOrderDao(ctx)
	order = &model.Order{
		UserID:    uId,
		ProductID: service.ProductId,
		BossID:    boss.ID,
		Num:       service.Num,
		Money:     float64(money) * float64(service.Num), // 计算出总价
		AddressId: service.AddressId,
		Type:      1,        // 默认未支付
		OrderNum:  orderNum, // 订单号
	}

	err = orderDao.CreateOrder(order)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreateOrder err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 更新商品热度
	cache.UpdateProductHeat(service.ProductId, "purchase")

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildOrder(order, product, boss, user, address),
	}
}

func (service *OrderService) List(ctx context.Context, uId uint) serializer.Response {
	code := e.Success

	if service.PageSize == 0 {
		service.PageSize = 15
	}

	condition := make(map[string]interface{})
	if service.Type != 0 { //
		condition["type"] = service.Type
	}
	condition["user_id"] = uId

	orderDao := dao.NewOrderDao(ctx)
	orderList, total, err := orderDao.ListOrderCondition(condition, service.BasePage)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ListOrder err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildOrders(ctx, orderList), uint(total))
}

func (service *OrderService) Show(ctx context.Context, uId uint, oId string) serializer.Response {
	code := e.Success
	orderDao := dao.NewOrderDao(ctx)
	orderId, _ := strconv.Atoi(oId)
	order, err := orderDao.ShowOrderByOid(uint(orderId), uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ListOrder err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	userDao := dao.NewUserDao(ctx)
	boss, _ := userDao.GetUserById(service.BossId)
	user, _ := userDao.GetUserById(uId)

	addressDao := dao.NewAddressDao(ctx)
	address, _ := addressDao.GetAddressByAid(service.AddressId, uId)

	productDao := dao.NewProductDao(ctx)
	product, _ := productDao.GetProductById(service.ProductId)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildOrder(order, product, boss, user, address),
	}
}

func (service *OrderService) Delete(ctx context.Context, uId uint, oId string) serializer.Response {
	orderId, _ := strconv.Atoi(oId)
	code := e.Success
	orderDao := dao.NewOrderDao(ctx)
	err := orderDao.DeleteOrderByOrderId(uint(orderId), uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("DeleteOrder err:", err)
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
