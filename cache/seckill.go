package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SeckillStockKey 秒杀库存 Key
func SeckillStockKey(activityID uint) string {
	return fmt.Sprintf("seckill:stock:%d", activityID)
}

// SeckillUserKey 用户购买资格 Key
func SeckillUserKey(activityID, userID uint) string {
	return fmt.Sprintf("seckill:user:%d:%d", activityID, userID)
}

// SeckillOrderQueueKey 异步下单队列 Key
func SeckillOrderQueueKey(activityID uint) string {
	return fmt.Sprintf("seckill:queue:%d", activityID)
}

// InitSeckillStock 初始化秒杀库存到 Redis
func InitSeckillStock(activityID uint, stock int) error {
	ctx := context.Background()
	key := SeckillStockKey(activityID)
	return RedisClient.Set(ctx, key, stock, 24*time.Hour).Err()
}

// ExecSeckillLua 执行秒杀 Lua 脚本
func ExecSeckillLua(activityID, userID uint) (int, error) {
	ctx := context.Background()

	// Lua 脚本：原子操作
	script := redis.NewScript(`
		local stockKey = KEYS[1]
		local userKey = KEYS[2]

		-- 1. 检查库存
		local stock = tonumber(redis.call('GET', stockKey))
		if stock == nil or stock <= 0 then
			return -1  -- 库存不足
		end

		-- 2. 检查用户是否已购买
		local bought = redis.call('EXISTS', userKey)
		if bought == 1 then
			return -2  -- 已购买
		end

		-- 3. 扣减库存
		redis.call('DECR', stockKey)

		-- 4. 标记用户已购买
		redis.call('SET', userKey, '1', 'EX', 86400)

		return 1  -- 成功
	`)

	keys := []string{
		SeckillStockKey(activityID),
		SeckillUserKey(activityID, userID),
	}

	result, err := script.Run(ctx, RedisClient, keys).Int()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// SeckillOrderInfo 秒杀订单信息
type SeckillOrderInfo struct {
	UserID     uint  `json:"user_id"`
	ActivityID uint  `json:"activity_id"`
	Timestamp  int64 `json:"timestamp"`
}

// PushSeckillOrder 推送秒杀订单到队列
func PushSeckillOrder(activityID uint, orderInfo *SeckillOrderInfo) error {
	ctx := context.Background()
	key := SeckillOrderQueueKey(activityID)

	data, err := json.Marshal(orderInfo)
	if err != nil {
		return err
	}

	return RedisClient.RPush(ctx, key, data).Err()
}

// PopSeckillOrder 从队列弹出秒杀订单
func PopSeckillOrder(activityID uint) (*SeckillOrderInfo, error) {
	ctx := context.Background()
	key := SeckillOrderQueueKey(activityID)

	result, err := RedisClient.LPop(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 队列为空
	}
	if err != nil {
		return nil, err
	}

	var orderInfo SeckillOrderInfo
	if err := json.Unmarshal([]byte(result), &orderInfo); err != nil {
		return nil, err
	}

	return &orderInfo, nil
}

// GetSeckillStock 获取秒杀库存
func GetSeckillStock(activityID uint) (int, error) {
	ctx := context.Background()
	key := SeckillStockKey(activityID)

	result, err := RedisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return result, nil
}

// HasUserBought 检查用户是否已购买
func HasUserBought(activityID, userID uint) (bool, error) {
	ctx := context.Background()
	key := SeckillUserKey(activityID, userID)

	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
