package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	AISessionTTL    = 24 * time.Hour // 会话过期时间
	MaxHistoryCount = 20             // 最大保留对话轮数
)

// AISessionContext 会话上下文结构
type AISessionContext struct {
	SessionID string            `json:"session_id"`
	UserID    uint              `json:"user_id"`
	History   []AIMessage       `json:"history"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// AIMessage 单条对话消息
type AIMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// AISessionKey 会话存储 Key
func AISessionKey(sessionID string) string {
	return fmt.Sprintf("ai:session:%s", sessionID)
}

// AISessionListKey 用户会话列表 Key
func AISessionListKey(userID uint) string {
	return fmt.Sprintf("ai:sessions:user:%d", userID)
}

// SaveAISession 保存会话上下文
func SaveAISession(sessionID string, ctx *AISessionContext) error {
	key := AISessionKey(sessionID)

	data, err := json.Marshal(ctx)
	if err != nil {
		return err
	}

	// 使用 Hash 存储
	err = RedisClient.HSet(context.Background(), key,
		"context", string(data),
		"updated_at", time.Now().Unix(),
	).Err()
	if err != nil {
		return err
	}

	// 设置 TTL
	RedisClient.Expire(context.Background(), key, AISessionTTL)

	// 更新用户会话列表
	userSessionsKey := AISessionListKey(ctx.UserID)
	RedisClient.ZAdd(context.Background(), userSessionsKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: sessionID,
	})

	return nil
}

// GetAISession 获取会话上下文
func GetAISession(sessionID string) (*AISessionContext, error) {
	key := AISessionKey(sessionID)

	data, err := RedisClient.HGet(context.Background(), key, "context").Result()
	if err == redis.Nil {
		return nil, nil // 会话不存在
	}
	if err != nil {
		return nil, err
	}

	var ctx AISessionContext
	if err := json.Unmarshal([]byte(data), &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// AppendAIMessage 添加对话记录
func AppendAIMessage(sessionID string, userID uint, msg AIMessage) error {
	ctx, err := GetAISession(sessionID)
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = &AISessionContext{
			SessionID: sessionID,
			UserID:    userID,
			History:   []AIMessage{},
			Metadata:  make(map[string]string),
			CreatedAt: time.Now(),
		}
	}

	ctx.History = append(ctx.History, msg)
	ctx.UpdatedAt = time.Now()

	// 只保留最近 N 轮对话
	if len(ctx.History) > MaxHistoryCount*2 {
		ctx.History = ctx.History[len(ctx.History)-MaxHistoryCount*2:]
	}

	return SaveAISession(sessionID, ctx)
}

// GetUserSessions 获取用户最近的会话列表
func GetUserSessions(userID uint, limit int64) ([]string, error) {
	key := AISessionListKey(userID)
	return RedisClient.ZRevRange(context.Background(), key, 0, limit-1).Result()
}

// DeleteAISession 删除会话
func DeleteAISession(sessionID string, userID uint) error {
	key := AISessionKey(sessionID)
	userSessionsKey := AISessionListKey(userID)

	RedisClient.Del(context.Background(), key)
	RedisClient.ZRem(context.Background(), userSessionsKey, sessionID)

	return nil
}
