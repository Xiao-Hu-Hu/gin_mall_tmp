package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// GetPayStatus 查询支付状态
func GetPayStatus(c *gin.Context) {
	orderIdStr := c.Param("id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "订单 ID 无效",
		})
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

	c.JSON(http.StatusOK, service.GetPayStatus(c.Request.Context(), uint(orderId), claims.ID))
}
