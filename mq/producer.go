package mq

import (
	"context"
	"encoding/json"
	"gin_mall_tmp/pkg/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishPayMessage 发送支付消息
func PublishPayMessage(msg *PayMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Channel.PublishWithContext(
		context.Background(),
		ExchangeName,  // exchange
		RoutingKeyPay, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 持久化
			Body:         body,
		},
	)
	if err != nil {
		util.LogrusObj.Infoln("发送支付消息失败:", err)
		return err
	}

	util.LogrusObj.Infoln("支付消息已发送:", msg.OrderNum)
	return nil
}

// PublishInventoryMessage 发送库存扣减消息
func PublishInventoryMessage(msg *InventoryMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Channel.PublishWithContext(
		context.Background(),
		ExchangeName,        // exchange
		RoutingKeyInventory, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		util.LogrusObj.Infoln("发送库存消息失败:", err)
		return err
	}

	util.LogrusObj.Infoln("库存扣减消息已发送: OrderID", msg.OrderID)
	return nil
}

// PublishCancelMessage 发送取消订单消息
func PublishCancelMessage(msg *OrderCancelMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Channel.PublishWithContext(
		context.Background(),
		ExchangeName,
		RoutingKeyCancel,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		util.LogrusObj.Infoln("发送取消消息失败:", err)
		return err
	}

	util.LogrusObj.Infoln("取消订单消息已发送: OrderID", msg.OrderID)
	return nil
}
