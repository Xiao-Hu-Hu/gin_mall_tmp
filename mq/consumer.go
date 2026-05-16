package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/pkg/util"
	"strconv"
)

// StartPayConsumer 启动支付消费者
func StartPayConsumer() {
	msgs, err := Channel.Consume(
		OrderPayQueue, // queue
		"",            // consumer
		false,         // auto-ack (手动确认)
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("注册支付消费者失败: %v", err))
	}

	go func() {
		for msg := range msgs {
			var payMsg PayMessage
			if err := json.Unmarshal(msg.Body, &payMsg); err != nil {
				util.LogrusObj.Infoln("解析支付消息失败:", err)
				msg.Nack(false, true) // 重新入队
				continue
			}

			if err := processPayMessage(&payMsg); err != nil {
				util.LogrusObj.Infoln("处理支付消息失败:", err)
				msg.Nack(false, true) // 重新入队
				continue
			}

			msg.Ack(false) // 确认消息
		}
	}()

	util.LogrusObj.Info("支付消费者已启动")
}

// processPayMessage 处理支付消息
func processPayMessage(msg *PayMessage) error {
	ctx := context.Background()
	orderDao := dao.NewOrderDao(ctx)
	userDao := dao.NewUserDao(ctx)

	// 1. 查询订单
	order, err := orderDao.ShowOrderByOid(msg.OrderID, msg.UserID)
	if err != nil {
		return fmt.Errorf("查询订单失败: %v", err)
	}

	// 2. 检查订单状态（防止重复支付）
	if order.Type == 2 {
		util.LogrusObj.Infoln("订单已支付，跳过:", msg.OrderNum)
		return nil
	}

	// 3. 开启事务
	tx := orderDao.Begin()

	// 4. 扣减用户余额
	user, err := userDao.GetUserById(msg.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询用户失败: %v", err)
	}

	moneyStr := util.Encrypt.AesDecoding(user.Money)
	moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)

	if moneyFloat < msg.Money {
		tx.Rollback()
		// 发送取消消息
		PublishCancelMessage(&OrderCancelMessage{
			OrderID: msg.OrderID,
			UserID:  msg.UserID,
		})
		return fmt.Errorf("余额不足")
	}

	// 扣减余额
	newMoney := fmt.Sprintf("%f", moneyFloat-msg.Money)
	user.Money = util.Encrypt.AesEncoding(newMoney)
	if err := userDao.UpdateUserById(msg.UserID, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新用户余额失败: %v", err)
	}

	// 5. 增加商家余额
	boss, err := userDao.GetUserById(msg.BossID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询商家失败: %v", err)
	}

	bossMoneyStr := util.Encrypt.AesDecoding(boss.Money)
	bossMoneyFloat, _ := strconv.ParseFloat(bossMoneyStr, 64)
	newBossMoney := fmt.Sprintf("%f", bossMoneyFloat+msg.Money)
	boss.Money = util.Encrypt.AesEncoding(newBossMoney)
	if err := userDao.UpdateUserById(msg.BossID, boss); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新商家余额失败: %v", err)
	}

	// 6. 更新订单状态
	order.Type = 2 // 已支付
	if err := orderDao.UpdateOrderByOrderId(msg.OrderID, order, msg.UserID); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 8. 发送库存扣减消息
	PublishInventoryMessage(&InventoryMessage{
		OrderID:   msg.OrderID,
		ProductID: msg.ProductID,
		Num:       msg.Num,
	})

	util.LogrusObj.Infoln("支付处理成功:", msg.OrderNum)
	return nil
}

// StartInventoryConsumer 启动库存消费者
func StartInventoryConsumer() {
	msgs, err := Channel.Consume(
		InventoryQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("注册库存消费者失败: %v", err))
	}

	go func() {
		for msg := range msgs {
			var invMsg InventoryMessage
			if err := json.Unmarshal(msg.Body, &invMsg); err != nil {
				util.LogrusObj.Infoln("解析库存消息失败:", err)
				msg.Nack(false, true)
				continue
			}

			if err := processInventoryMessage(&invMsg); err != nil {
				util.LogrusObj.Infoln("处理库存消息失败:", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	util.LogrusObj.Info("库存消费者已启动")
}

// processInventoryMessage 处理库存扣减
func processInventoryMessage(msg *InventoryMessage) error {
	ctx := context.Background()
	productDao := dao.NewProductDao(ctx)

	product, err := productDao.GetProductById(msg.ProductID)
	if err != nil {
		return fmt.Errorf("查询商品失败: %v", err)
	}

	if product.Num < msg.Num {
		return fmt.Errorf("库存不足: 当前 %d, 需要 %d", product.Num, msg.Num)
	}

	product.Num -= msg.Num
	if err := productDao.UpdateProductById(msg.ProductID, product); err != nil {
		return fmt.Errorf("更新库存失败: %v", err)
	}

	util.LogrusObj.Infoln("库存扣减成功: ProductID", msg.ProductID)
	return nil
}

// StartCancelConsumer 启动取消订单消费者
func StartCancelConsumer() {
	msgs, err := Channel.Consume(
		OrderCancelQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("注册取消订单消费者失败: %v", err))
	}

	go func() {
		for msg := range msgs {
			var cancelMsg OrderCancelMessage
			if err := json.Unmarshal(msg.Body, &cancelMsg); err != nil {
				util.LogrusObj.Infoln("解析取消消息失败:", err)
				msg.Nack(false, true)
				continue
			}

			if err := processCancelMessage(&cancelMsg); err != nil {
				util.LogrusObj.Infoln("处理取消消息失败:", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	util.LogrusObj.Info("取消订单消费者已启动")
}

// processCancelMessage 处理取消订单
func processCancelMessage(msg *OrderCancelMessage) error {
	ctx := context.Background()
	orderDao := dao.NewOrderDao(ctx)

	// 更新订单状态为已取消 (Type=3)
	err := orderDao.UpdateOrderByOrderId(msg.OrderID, nil, msg.UserID)
	if err != nil {
		return fmt.Errorf("取消订单失败: %v", err)
	}

	util.LogrusObj.Infoln("订单已取消: OrderID", msg.OrderID)
	return nil
}
