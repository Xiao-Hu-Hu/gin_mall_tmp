package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 创建商品
func CreateFavorite(c *gin.Context) {

	creatFavoriteService := service.FavoriteService{}

	// 先绑定参数
	if err := c.ShouldBind(&creatFavoriteService); err != nil {
		util.LogrusObj.Infoln("creatFavorite err: ", err)
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
		util.LogrusObj.Infoln("creatFavorite err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := creatFavoriteService.Create(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

func ShowFavorite(c *gin.Context) {
	var showFavoriteService service.FavoriteService

	// 先绑定参数
	if err := c.ShouldBind(&showFavoriteService); err != nil {
		util.LogrusObj.Infoln("creatFavorite err: ", err)
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
		util.LogrusObj.Infoln("showFavorite err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := showFavoriteService.Show(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

func DeleteFavorite(c *gin.Context) {
	var deleteFavoriteService service.FavoriteService

	// 先绑定参数
	if err := c.ShouldBind(&deleteFavoriteService); err != nil {
		util.LogrusObj.Infoln("deleteFavorite err: ", err)
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
		util.LogrusObj.Infoln("deleteFavorite err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := deleteFavoriteService.Delete(c.Request.Context(), claims.ID, c.Param("id"))
	c.JSON(http.StatusOK, res)
	return
}
