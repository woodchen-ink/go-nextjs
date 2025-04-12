package router

import (
	"go-nextjs/handler"
	"go-nextjs/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 认证相关路由
	auth := r.Group("/api/auth")
	{
		auth.GET("/login", handler.Login)
		auth.GET("/callback", handler.Callback)
		auth.GET("/me", middleware.AuthRequired(), handler.GetCurrentUser)
		auth.GET("/logout", handler.Logout)
	}

	// 公开路由
	public := r.Group("/api")
	{
		//返回ok
		public.GET("/deals", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})
	}

	// 需要认证的路由 - 只需要验证用户身份，不再检查admin角色
	admin := r.Group("/api")
	admin.Use(middleware.AuthRequired())
	{

	}

	return r
}
