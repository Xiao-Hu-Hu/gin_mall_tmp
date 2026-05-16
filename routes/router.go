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
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		// 轮播图
		v1.GET("carousels", api.ListCarousel)

		// 商品操作
		v1.GET("products", api.ListProduct)
		v1.GET("products/:id", api.ShowProduct)
		v1.GET("imgs/:id", api.ListProductImg)
		v1.GET("categories", api.ListCategory)

		// 排行榜（公开接口）
		v1.GET("products/hot", api.GetHotProducts)
		v1.GET("products/view-rank", api.GetViewRankProducts)
		v1.GET("products/purchase-rank", api.GetPurchaseRankProducts)
		v1.GET("products/favorite-rank", api.GetFavoriteRankProducts)

		// 秒杀活动（公开接口）
		v1.GET("seckill/:id", api.GetSeckillActivity)

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.PUT("user", api.UserUpdate)
			authed.POST("avatar", api.UploadAvatar)
			authed.POST("user/sending-email", api.SendEmail)
			authed.POST("user/valid-email", api.ValidEmail)

			authed.POST("money", api.ShowMoney)

			authed.POST("product", api.CreateProduct)
			authed.POST("products", api.SearchProduct)

			authed.GET("favorites", api.ShowFavorite)
			authed.POST("favorites", api.CreateFavorite)
			authed.DELETE("favorites/:id", api.DeleteFavorite)

			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.GetAddress)
			authed.GET("addresses", api.ListAddress)
			authed.PUT("addresses/:id", api.UpdateAddress)
			authed.DELETE("addresses/:id", api.DeleteAddress)

			authed.POST("carts", api.CreateCart)
			authed.GET("carts", api.ListCart)
			authed.GET("carts/:id", api.ShowCart)
			authed.PUT("carts/:id", api.UpdateCart)
			authed.DELETE("carts/:id", api.DeleteCart)

			authed.POST("orders", api.CreateOrder)
			authed.GET("orders", api.ListOrder)
			authed.GET("orders/:id", api.ShowOrder)
			authed.DELETE("orders/:id", api.DeleteOrder)

			authed.POST("ai/customer-service/chat", api.CustomerServiceChat)
			authed.GET("ai/sessions", api.GetAISessions)
			authed.GET("ai/sessions/:session_id", api.GetAISessionHistory)
			authed.DELETE("ai/sessions/:session_id", api.DeleteAISession)

			// RAG 知识库管理
			authed.POST("ai/knowledge/sync", api.SyncKnowledge)
			authed.POST("ai/knowledge/search", api.SearchKnowledge)
			authed.GET("ai/knowledge/stats", api.GetKnowledgeStats)

			// 秒杀订单（需认证）
			authed.POST("seckill/order", api.CreateSeckillOrder)
			authed.GET("seckill/order/:activity_id", api.CheckSeckillOrder)

			authed.POST("paydown", api.OrderPay)
			authed.GET("pay/status/:id", api.GetPayStatus)
		}
	}

	return r
}
