package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateAddress(c *gin.Context) {

	creatAddressService := service.AddressService{}

	// 先绑定参数
	if err := c.ShouldBind(&creatAddressService); err != nil {
		util.LogrusObj.Infoln("creatAddress err: ", err)
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
		util.LogrusObj.Infoln("creatAddress err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := creatAddressService.Create(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

func GetAddress(c *gin.Context) {
	getAddressService := service.AddressService{}

	// 先绑定参数
	if err := c.ShouldBind(&getAddressService); err != nil {
		util.LogrusObj.Infoln("getAddress err: ", err)
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
		util.LogrusObj.Infoln("getAddress err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := getAddressService.Get(c.Request.Context(), claims.ID, c.Param("id"))
	c.JSON(http.StatusOK, res)
	return
}

func ListAddress(c *gin.Context) {
	var listAddressService service.AddressService

	// 先绑定参数
	if err := c.ShouldBind(&listAddressService); err != nil {
		util.LogrusObj.Infoln("ListAddress err: ", err)
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
		util.LogrusObj.Infoln("listAddress err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := listAddressService.List(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

func UpdateAddress(c *gin.Context) {

	updateAddressService := service.AddressService{}

	// 先绑定参数
	if err := c.ShouldBind(&updateAddressService); err != nil {
		util.LogrusObj.Infoln("updateAddress err: ", err)
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
		util.LogrusObj.Infoln("updateAddress err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := updateAddressService.Update(c.Request.Context(), claims.ID, c.Param("id"))
	c.JSON(http.StatusOK, res)
	return
}

func DeleteAddress(c *gin.Context) {
	var deleteAddressService service.AddressService

	// 先绑定参数
	if err := c.ShouldBind(&deleteAddressService); err != nil {
		util.LogrusObj.Infoln("deleteAddress err: ", err)
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
		util.LogrusObj.Infoln("deleteAddress err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	res := deleteAddressService.Delete(c.Request.Context(), claims.ID, c.Param("id"))
	c.JSON(http.StatusOK, res)
	return
}
