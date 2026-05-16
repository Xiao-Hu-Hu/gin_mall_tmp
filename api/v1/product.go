package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 创建商品
func CreateProduct(c *gin.Context) {

	creatProductService := service.ProductService{}

	// 先绑定参数
	if err := c.ShouldBind(&creatProductService); err != nil {
		util.LogrusObj.Infoln("creatProduct err: ", err)
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
		util.LogrusObj.Infoln("creatProduct err: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌验证失败",
			"error":  err.Error(),
		})
		return
	}

	// 处理文件上传
	form, _ := c.MultipartForm()
	files := form.File["file"] // 获取所有上传的商品图片文件

	res := creatProductService.Create(c.Request.Context(), claims.ID, files)
	c.JSON(http.StatusOK, res)
	return
}

func ListProduct(c *gin.Context) {
	var listProductService service.ProductService

	if err := c.ShouldBind(&listProductService); err == nil {
		res := listProductService.List(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln("ListProduct err: ", err)
	}
}

func ShowProduct(c *gin.Context) {
	var showProductService service.ProductService

	if err := c.ShouldBind(&showProductService); err == nil {
		res := showProductService.Show(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln("showProduct err: ", err)
	}
}

func SearchProduct(c *gin.Context) {
	var searchProductService service.ProductService

	if err := c.ShouldBind(&searchProductService); err == nil {
		res := searchProductService.Search(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln("SearchProduct err: ", err)
	}
}
