package main

import (
	"log"
	"os"

	"go-nextjs/config"
	"go-nextjs/cron"
	"go-nextjs/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 打印运行环境信息
	log.Printf("数据目录: %s", os.Getenv("DATA_DIR"))

	// 初始化配置和数据库
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 确保数据库连接正常
	if config.DB == nil {
		log.Fatalf("数据库连接未正确初始化")
	}

	// 初始化定时任务
	if err := cron.Init(); err != nil {
		log.Fatalf("初始化定时任务失败: %v", err)
	}

	// 设置gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
