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
	"time"
)

type SeckillService struct {
	ActivityID uint `json:"activity_id" form:"activity_id"`
}

// CreateSeckillOrder 创建秒杀订单
func (service *SeckillService) CreateSeckillOrder(ctx context.Context, userID uint) serializer.Response {
	code := e.Success

	// 1. 参数校验
	if service.ActivityID == 0 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "活动 ID 不能为空",
		}
	}

	// 2. 检查活动是否存在且进行中
	seckillDao := dao.NewSeckillDao(ctx)
	activity, err := seckillDao.GetActiveActivityByID(service.ActivityID)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("GetActiveActivityByID err:", err)
		return serializer.Response{
			Status: code,
			Msg:    "活动不存在或未开始",
			Error:  err.Error(),
		}
	}

	// 3. 检查活动时间
	now := time.Now()
	if now.Before(activity.StartAt) {
		return serializer.Response{
			Status: code,
			Msg:    "秒杀尚未开始",
		}
	}
	if now.After(activity.EndAt) {
		return serializer.Response{
			Status: code,
			Msg:    "秒杀已结束",
		}
	}

	// 4. 执行 Lua 脚本（原子操作：检查库存 + 判断资格 + 扣减库存）
	result, err := cache.ExecSeckillLua(service.ActivityID, userID)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ExecSeckillLua err:", err)
		return serializer.Response{
			Status: code,
			Msg:    "秒杀失败",
			Error:  err.Error(),
		}
	}

	if result == -1 {
		return serializer.Response{
			Status: code,
			Msg:    "库存不足",
		}
	}
	if result == -2 {
		return serializer.Response{
			Status: code,
			Msg:    "您已购买过",
		}
	}

	// 5. 发送到异步下单队列
	orderInfo := &cache.SeckillOrderInfo{
		UserID:     userID,
		ActivityID: service.ActivityID,
		Timestamp:  now.Unix(),
	}
	err = cache.PushSeckillOrder(service.ActivityID, orderInfo)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("PushSeckillOrder err:", err)
		return serializer.Response{
			Status: code,
			Msg:    "下单失败",
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    "排队中，请稍后查询订单",
	}
}

// ProcessSeckillOrders 处理异步下单队列
func ProcessSeckillOrders() {
	for {
		// 这里需要遍历所有活动，实际项目中可以使用配置或数据库查询
		// 简化实现：假设只有一个活动
		activityID := uint(1)

		orderInfo, err := cache.PopSeckillOrder(activityID)
		if err != nil {
			util.LogrusObj.Infoln("PopSeckillOrder err:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if orderInfo == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 创建订单
		processSeckillOrder(orderInfo)
	}
}

func processSeckillOrder(orderInfo *cache.SeckillOrderInfo) {
	ctx := context.Background()
	seckillDao := dao.NewSeckillDao(ctx)

	// 获取活动信息
	activity, err := seckillDao.GetActivityByID(orderInfo.ActivityID)
	if err != nil {
		util.LogrusObj.Infoln("GetActivityByID err:", err)
		return
	}

	// 生成订单号
	orderNum := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000000))
	orderNum = fmt.Sprintf("SK%s%d%d", orderNum, activity.ProductID, orderInfo.UserID)

	// 创建秒杀订单
	order := &model.SeckillOrder{
		UserID:     orderInfo.UserID,
		ActivityID: orderInfo.ActivityID,
		ProductID:  activity.ProductID,
		OrderNum:   orderNum,
		Money:      activity.SeckillPrice,
		Status:     0, // 待支付
	}

	err = seckillDao.CreateSeckillOrder(order)
	if err != nil {
		util.LogrusObj.Infoln("CreateSeckillOrder err:", err)
		return
	}

	util.LogrusObj.Infoln("秒杀订单创建成功:", orderNum)
}

// GetSeckillActivity 获取秒杀活动详情
func GetSeckillActivity(ctx context.Context, activityID uint) serializer.Response {
	code := e.Success

	seckillDao := dao.NewSeckillDao(ctx)
	activity, err := seckillDao.GetActivityByID(activityID)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "获取活动失败",
			Error:  err.Error(),
		}
	}

	// 获取库存
	stock, _ := cache.GetSeckillStock(activityID)

	return serializer.Response{
		Status: code,
		Msg:    "ok",
		Data: map[string]interface{}{
			"activity":      activity,
			"remaining_stock": stock,
		},
	}
}

// CheckSeckillOrder 检查秒杀订单状态
func CheckSeckillOrder(ctx context.Context, userID, activityID uint) serializer.Response {
	code := e.Success

	seckillDao := dao.NewSeckillDao(ctx)
	order, err := seckillDao.GetSeckillOrderByUserAndActivity(userID, activityID)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "未找到订单",
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    "ok",
		Data:   order,
	}
}
