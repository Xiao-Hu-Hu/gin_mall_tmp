package v1

import (
	"gin_mall_tmp/pkg/util"
	aiservice "gin_mall_tmp/service/ai"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CustomerServiceChat(c *gin.Context) {
	var request aiservice.CustomerServiceRequest

	contentType := c.ContentType()
	if strings.Contains(contentType, "application/json") {
		// JSON 格式
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
	} else {
		// form-data 或 x-www-form-urlencoded 格式
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
	}

	token := util.ExtractTokenFromRequest(c.Request)
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
			"msg":    "令牌校验失败",
			"error":  err.Error(),
		})
		return
	}

	service := aiservice.NewCustomerService(&request, claims.ID, token)
	c.JSON(http.StatusOK, service.Chat(c.Request.Context()))
}
