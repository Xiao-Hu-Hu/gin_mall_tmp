package cache

import (
	"context"
	"fmt"
	"gin_mall_tmp/pkg/util"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPwd    string
	RedisDbName string
)

func init() {
	var file *ini.File
	var err error

	// 优先尝试加载Docker配置文件
	file, err = ini.Load("./conf/config.docker.ini")
	if err != nil {
		// 如果Docker配置文件不存在，则加载默认配置
		file, err = ini.Load("./conf/config.ini")
		if err != nil {
			fmt.Println("redis config err", err)
			panic(err)
		}
	}
	LoadRedisData(file)
	Redis()
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPwd = file.Section("redis").Key("RedisPassword").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPwd,
		DB:       int(db),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		util.LogrusObj.Infoln("redis err: ", err)
		// 在开发环境下，如果Redis连接失败，不panic，而是记录警告
		// 这样应用可以继续运行，但验证码功能将不可用
		fmt.Printf("警告: Redis连接失败，验证码功能将不可用: %v\n", err)
		RedisClient = nil
		return
	}
	RedisClient = client
}
