package v1

import (
	"gin_mall_tmp/cache"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/serializer"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetHotProducts 获取综合热度排行榜
func GetHotProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	productIDs, err := cache.GetHotProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "获取排行榜失败",
			"error":  err.Error(),
		})
		return
	}

	products := getProductsByIDs(productIDs)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "ok",
		"data":   serializer.BuildProducts(products),
	})
}

// GetViewRankProducts 获取浏览量排行榜
func GetViewRankProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	productIDs, err := cache.GetViewRankProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "获取排行榜失败",
			"error":  err.Error(),
		})
		return
	}

	products := getProductsByIDs(productIDs)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "ok",
		"data":   serializer.BuildProducts(products),
	})
}

// GetPurchaseRankProducts 获取购买量排行榜
func GetPurchaseRankProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	productIDs, err := cache.GetPurchaseRankProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "获取排行榜失败",
			"error":  err.Error(),
		})
		return
	}

	products := getProductsByIDs(productIDs)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "ok",
		"data":   serializer.BuildProducts(products),
	})
}

// GetFavoriteRankProducts 获取收藏量排行榜
func GetFavoriteRankProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	productIDs, err := cache.GetFavoriteRankProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "获取排行榜失败",
			"error":  err.Error(),
		})
		return
	}

	products := getProductsByIDs(productIDs)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "ok",
		"data":   serializer.BuildProducts(products),
	})
}

// getProductsByIDs 根据ID列表获取商品
func getProductsByIDs(productIDs []string) []*model.Product {
	if len(productIDs) == 0 {
		return nil
	}

	products := make([]*model.Product, 0, len(productIDs))
	for _, idStr := range productIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		product, err := dao.NewProductDao(nil).GetProductById(uint(id))
		if err != nil {
			continue
		}
		products = append(products, product)
	}
	return products
}
