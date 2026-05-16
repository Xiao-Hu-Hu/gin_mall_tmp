package main

import (
	"context"
	"log"
	"os"
	"time"

	"gin_mall_tmp/conf"
	"gin_mall_tmp/mq"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/routes"
	aiservice "gin_mall_tmp/service/ai"

	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		util.LogrusObj.Infoln("未找到 .env 文件，使用系统环境变量")
	}

	conf.Init()

	// 初始化 RabbitMQ
	mq.InitRabbitMQ(conf.RabbitMQ)
	defer mq.Close()

	// 启动消费者
	mq.StartPayConsumer()
	mq.StartInventoryConsumer()
	mq.StartCancelConsumer()

	// 初始化 Milvus 向量数据库
	initMilvus()

	// 启动 HTTP 服务
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}

func initMilvus() {
	milvusAddr := os.Getenv("MILVUS_ADDR")
	if milvusAddr == "" {
		milvusAddr = "localhost:19530"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{
		Address: milvusAddr,
	})
	if err != nil {
		util.LogrusObj.Warnf("Milvus 连接失败 (%s)，向量检索功能不可用: %v", milvusAddr, err)
		return
	}

	aiservice.SetMilvusClient(c)

	if err := aiservice.InitMilvusCollection(ctx); err != nil {
		util.LogrusObj.Warnf("Milvus collection 初始化失败: %v", err)
		return
	}

	// 确保 collection 已加载
	_ = aiservice.EnsureCollectionLoaded(ctx)

	log.Printf("[Milvus] connected to %s, collection ready", milvusAddr)
}
