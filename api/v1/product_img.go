package v1

import (
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListProductImg(c *gin.Context) {
	//用户注册
	var listProductImg service.ListProductImg

	if err := c.ShouldBind(&listProductImg); err == nil {
		res := listProductImg.List(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
