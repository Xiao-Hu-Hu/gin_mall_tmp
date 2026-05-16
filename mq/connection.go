package mq

import (
	"fmt"
	"gin_mall_tmp/pkg/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Conn    *amqp.Connection
	Channel *amqp.Channel
)

const (
	ExchangeName = "mall.exchange"

	OrderPayQueue    = "order.pay.queue"
	OrderCancelQueue = "order.cancel.queue"
	InventoryQueue   = "inventory.queue"

	RoutingKeyPay       = "order.pay"
	RoutingKeyCancel    = "order.cancel"
	RoutingKeyInventory = "inventory.deduct"
)

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(url string) {
	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		panic(fmt.Sprintf("连接 RabbitMQ 失败: %v", err))
	}

	Channel, err = Conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("创建 Channel 失败: %v", err))
	}

	// 声明 Exchange
	err = Channel.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		panic(fmt.Sprintf("声明 Exchange 失败: %v", err))
	}

	// 声明队列并绑定
	declareQueue(OrderPayQueue, RoutingKeyPay)
	declareQueue(OrderCancelQueue, RoutingKeyCancel)
	declareQueue(InventoryQueue, RoutingKeyInventory)

	util.LogrusObj.Info("RabbitMQ 初始化成功")
}

func declareQueue(name, routingKey string) {
	_, err := Channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(fmt.Sprintf("声明队列 %s 失败: %v", name, err))
	}

	err = Channel.QueueBind(
		name,         // queue name
		routingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("绑定队列 %s 失败: %v", name, err))
	}
}

// Close 关闭连接
func Close() {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
}
