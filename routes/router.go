package routes

import (
	api "gin_mall_tmp/api/v1"
	"gin_mall_tmp/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.StaticFS("/static", http.Dir("./static"))
	r.Static("/img", "./static/img")
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})

		//用户操作
		v1.POST("user/register", api.UserRegister)           //注册
		v1.POST("user/login", api.UserLogin)                 //登录
		v1.POST("user/send-verify-code", api.SendEmailVerifyCode) //发送邮箱验证码
		v1.POST("user/email-register", api.EmailRegister)        //邮箱验证码注册

		// 轮播图
		v1.GET("carousels", api.ListCarousel) //展示轮播图

		// 商品操作
		v1.GET("products", api.ListProduct)     //获取商品列表
		v1.GET("products/:id", api.ShowProduct) //展示商品详细信息
		v1.GET("imgs/:id", api.ListProductImg)  //获取商品图片地址
		v1.GET("categories", api.ListCategory)  //获取商品分类

		authed := v1.Group("/") // 需要登录保护
		authed.Use(middleware.JWT())
		{
			// 用户操作
			authed.PUT("user", api.UserUpdate)               //更新用户昵称
			authed.POST("avatar", api.UploadAvatar)          //上传用户头像
			authed.POST("user/sending-email", api.SendEmail) //发送邮件
			authed.POST("user/valid-email", api.ValidEmail)  //解析邮件

			// 显示金额
			authed.POST("money", api.ShowMoney) //显示用户余额

			// 商品操作
			authed.POST("product", api.CreateProduct)  //创建商品
			authed.POST("products", api.SearchProduct) //搜索商品

			// 收藏夹操作
			authed.GET("favorites", api.ShowFavorite)          //展示收藏夹列表
			authed.POST("favorites", api.CreateFavorite)       //创建收藏
			authed.DELETE("favorites/:id", api.DeleteFavorite) //取消收藏

			// 地址操作
			authed.POST("addresses", api.CreateAddress)       //创建地址
			authed.GET("addresses/:id", api.GetAddress)       //展示单个地址
			authed.GET("addresses", api.ListAddress)          //展示所有地址
			authed.PUT("addresses/:id", api.UpdateAddress)    //更新地址
			authed.DELETE("addresses/:id", api.DeleteAddress) //删除地址

			// 购物车操作
			authed.POST("carts", api.CreateCart)       //创建购物车
			authed.GET("carts", api.ListCart)          //展示所有购物车
			authed.GET("carts/:id", api.ShowCart)      //展示单个购物车
			authed.PUT("carts/:id", api.UpdateCart)    //修改购物车
			authed.DELETE("carts/:id", api.DeleteCart) //删除购物车

			// 订单操作
			authed.POST("orders", api.CreateOrder)       //创建订单
			authed.GET("orders", api.ListOrder)          //展示所有订单
			authed.GET("orders/:id", api.ShowOrder)      //展示单个订单
			authed.DELETE("orders/:id", api.DeleteOrder) //删除订单

			// 支付功能
			authed.POST("paydown", api.OrderPay)
		}
	}

	return r
}
