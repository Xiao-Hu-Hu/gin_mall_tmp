package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// user 路由处理
func UserRegister(c *gin.Context) {
	//用户注册
	var userRegister service.UserService

	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln("UserRegister err: ", err)
	}
}

func UserLogin(c *gin.Context) {
	//用户登录
	var userLogin service.UserService

	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln("UserLogin err: ", err)
	}
}

// 需要保护的
func UserUpdate(c *gin.Context) {
	//用户更新信息
	var userUpdate service.UserService

	// 先绑定参数
	if err := c.ShouldBind(&userUpdate); err != nil {
		util.LogrusObj.Infoln("UserUpdate err: ", err)
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
		util.LogrusObj.Infoln("UserUpdate err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	//更新信息
	res := userUpdate.Update(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

func UploadAvatar(c *gin.Context) {
	// 上传用户头像
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		util.LogrusObj.Infoln("UploadAvatar err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "文件上传失败",
			"error":  err.Error(),
		})
		return
	}

	if fileHeader == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "文件信息为空",
		})
		return
	}

	fileSize := fileHeader.Size
	var uploadAvatar service.UserService

	// 先绑定参数
	if err := c.ShouldBind(&uploadAvatar); err != nil {
		util.LogrusObj.Infoln("UploadAvatar err: ", err)
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
		util.LogrusObj.Infoln("UploadAvatar err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	// 上传图片
	res := uploadAvatar.Post(c.Request.Context(), claims.ID, file, fileSize)
	c.JSON(http.StatusOK, res)

}

// email 路由处理

// 需要保护的
func SendEmail(c *gin.Context) {
	// 发送邮件
	var sendEmail service.SendEmailService

	// 先绑定参数
	if err := c.ShouldBind(&sendEmail); err != nil {
		util.LogrusObj.Infoln("SendEmail err: ", err)
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
		util.LogrusObj.Infoln("SendEmail err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}
	//更新信息
	res := sendEmail.Send(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}

// ValidEmail
func ValidEmail(c *gin.Context) {
	// 解析邮件token
	var validEmail service.ValidEmailService

	// 先绑定参数
	if err := c.ShouldBind(&validEmail); err != nil {
		util.LogrusObj.Infoln("ValidEmail err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "绑定参数失败",
			"error":  err.Error(),
		})
		return
	}

	//更新信息
	res := validEmail.Valid(c.Request.Context(), c.GetHeader("Authorization"))
	c.JSON(http.StatusOK, res)
	return
}

// ShowMoney
func ShowMoney(c *gin.Context) {
	var showMoney service.ShowMoneyService

	// 先绑定参数
	if err := c.ShouldBind(&showMoney); err != nil {
		util.LogrusObj.Infoln("ShowMoney err: ", err)
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
		util.LogrusObj.Infoln("ShowMoney err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}
	//更新信息
	res := showMoney.Show(c.Request.Context(), claims.ID)
	c.JSON(http.StatusOK, res)
	return
}
