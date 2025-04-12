package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Init 初始化配置
func Init() error {
	// 初始化环境变量
	if err := LoadEnv(); err != nil {
		return err
	}

	// 初始化数据库
	if err := InitDB(); err != nil {
		return err
	}

	return nil
}

// LoadEnv 加载环境变量
func LoadEnv() error {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 设置默认值
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	// 加载所有配置
	if err := LoadConfig(); err != nil {
		return err
	}

	return nil
}
