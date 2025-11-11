package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func OrderPay(c *gin.Context) {
	var orderPayService service.OrderPayService

	// 先绑定参数
	if err := c.ShouldBind(&orderPayService); err != nil {
		util.LogrusObj.Infoln("orderPay err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "绑定参数失败",
			"error":  err.Error(),
		})
		return
	}

	// 获取并验证token
	token := c.GetHeader("token")
	if token == "" {
		// 如果 token 头为空，尝试用 Authorization 头
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
	// 解析token
	claims, err := util.ParseToken(token)
	if err != nil {
		util.LogrusObj.Infoln("PayService err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := orderPayService.PayDown(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}
