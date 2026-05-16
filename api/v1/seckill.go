package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// CreateSeckillOrder 创建秒杀订单
func CreateSeckillOrder(c *gin.Context) {
	var seckillService service.SeckillService

	if err := c.ShouldBind(&seckillService); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	token := c.GetHeader("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "未提供令牌",
		})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, seckillService.CreateSeckillOrder(c.Request.Context(), claims.ID))
}

// GetSeckillActivity 获取秒杀活动详情
func GetSeckillActivity(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "活动 ID 无效",
		})
		return
	}

	c.JSON(http.StatusOK, service.GetSeckillActivity(c.Request.Context(), uint(activityID)))
}

// CheckSeckillOrder 检查秒杀订单状态
func CheckSeckillOrder(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "未提供令牌",
		})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	activityIDStr := c.Param("activity_id")
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "活动 ID 无效",
		})
		return
	}

	c.JSON(http.StatusOK, service.CheckSeckillOrder(c.Request.Context(), claims.ID, uint(activityID)))
}
