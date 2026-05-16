package cache

import (
	"context"
	"strconv"
)

// UpdateProductHeat 更新商品热度
func UpdateProductHeat(productID uint, action string) {
	ctx := context.Background()
	pid := strconv.Itoa(int(productID))

	switch action {
	case "view":
		RedisClient.ZIncrBy(ctx, RankKey, 1, pid)
		RedisClient.ZIncrBy(ctx, RankHotKey, 1, pid)
	case "purchase":
		RedisClient.ZIncrBy(ctx, RankPurchaseKey, 1, pid)
		RedisClient.ZIncrBy(ctx, RankHotKey, 5, pid)
	case "favorite":
		RedisClient.ZIncrBy(ctx, RankFavoriteKey, 1, pid)
		RedisClient.ZIncrBy(ctx, RankHotKey, 3, pid)
	}
}

// GetHotProducts 获取综合热度 Top N
func GetHotProducts(topN int64) ([]string, error) {
	ctx := context.Background()
	return RedisClient.ZRevRange(ctx, RankHotKey, 0, topN-1).Result()
}

// GetViewRankProducts 获取浏览量 Top N
func GetViewRankProducts(topN int64) ([]string, error) {
	ctx := context.Background()
	return RedisClient.ZRevRange(ctx, RankKey, 0, topN-1).Result()
}

// GetPurchaseRankProducts 获取购买量 Top N
func GetPurchaseRankProducts(topN int64) ([]string, error) {
	ctx := context.Background()
	return RedisClient.ZRevRange(ctx, RankPurchaseKey, 0, topN-1).Result()
}

// GetFavoriteRankProducts 获取收藏量 Top N
func GetFavoriteRankProducts(topN int64) ([]string, error) {
	ctx := context.Background()
	return RedisClient.ZRevRange(ctx, RankFavoriteKey, 0, topN-1).Result()
}

// GetProductScore 获取商品在指定排行榜的分数
func GetProductScore(key string, productID uint) (float64, error) {
	ctx := context.Background()
	return RedisClient.ZScore(ctx, key, strconv.Itoa(int(productID))).Result()
}
