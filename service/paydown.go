package service

import (
	"context"
	"fmt"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/mq"
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
	code := e.Success

	// 1. 查询订单
	orderDao := dao.NewOrderDao(ctx)
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

	// 2. 检查订单状态
	if order.Type == 2 {
		return serializer.Response{
			Status: code,
			Msg:    "订单已支付",
		}
	}

	// 3. 校验余额
	userDao := dao.NewUserDao(ctx)
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

	util.Encrypt.SetKey(service.Key)
	moneyStr := util.Encrypt.AesDecoding(user.Money)
	moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)

	if moneyFloat < order.Money {
		return serializer.Response{
			Status: code,
			Msg:    "余额不足",
		}
	}

	// 4. 发送支付消息到 RabbitMQ
	payMsg := &mq.PayMessage{
		OrderID:   service.OrderId,
		UserID:    uId,
		ProductID: order.ProductID,
		BossID:    order.BossID,
		Money:     order.Money,
		OrderNum:  order.OrderNum,
		Num:       order.Num,
	}

	if err := mq.PublishPayMessage(payMsg); err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "发送支付请求失败",
			Error:  err.Error(),
		}
	}

	// 5. 立即返回
	return serializer.Response{
		Status: code,
		Msg:    "支付处理中，请稍后查看订单状态",
	}
}

// GetPayStatus 查询支付状态
func GetPayStatus(ctx context.Context, orderId uint, uId uint) serializer.Response {
	code := e.Success
	orderDao := dao.NewOrderDao(ctx)

	order, err := orderDao.ShowOrderByOid(orderId, uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	status := "待支付"
	if order.Type == 2 {
		status = "已支付"
	} else if order.Type == 3 {
		status = "已取消"
	}

	return serializer.Response{
		Status: code,
		Msg:    "ok",
		Data: map[string]interface{}{
			"order_id": order.ID,
			"order_num": order.OrderNum,
			"status":    status,
			"type":      order.Type,
			"money":     fmt.Sprintf("%.2f", order.Money),
		},
	}
}

// 需要补充 model.Order 的 Type 常量
const (
	OrderTypeUnpaid  = 1 // 未支付
	OrderTypePaid    = 2 // 已支付
	OrderTypeCancel  = 3 // 已取消
)

// 确保 model 包被引用
var _ = model.Order{}
