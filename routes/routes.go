package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// func CorsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 100))

	// 配置 CORS（开发环境允许所有来源）
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"}, // 前端地址
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	AllowCredentials: true, // 允许传 Cookie
	// 	MaxAge:           12 * time.Hour,
	// }))
	// r.Use(CorsMiddleware())
	// 或者最简单的允许所有来源（仅限开发测试）
	//r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetStringSlice("cors.allow_origins"),
		AllowMethods:     viper.GetStringSlice("cors.allow_methods"),
		AllowHeaders:     viper.GetStringSlice("cors.allow_headers"),
		AllowCredentials: viper.GetBool("cors.allow_credentials"),
		MaxAge:           time.Duration(viper.GetInt("cors.max_age")) * time.Hour,
	}))
	v1 := r.Group("/api/v1")
	//注册路由业务
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	// v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	c.Request.Header.Get("Authorization")
	// })
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.DELETE("/post/:id", controller.DeletePostHandler)
		v1.GET("posts/", controller.GetPostListHandler)
		v1.GET("posts2/", controller.GetPostListHandler2)
		v1.POST("uploadimg", controller.UploadImageHandler)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
