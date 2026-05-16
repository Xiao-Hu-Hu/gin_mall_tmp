package middleware

import (
	"fmt"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.Success
		token := c.Request.Header.Get("token")
		if token == "" {
			code = 404
			fmt.Println("middleware JWT: Token is empty")
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthToken
				fmt.Println("token err:", err)
			} else if time.Now().After(claims.ExpiresAt.Time) {
				code = e.ErrorAuthCheckTokenTimeOut
			}
		}
		if code != e.Success {
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
