package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	// AppEnv 应用环境
	AppEnv string
	// ServerPort 服务端口
	ServerPort string
	// DatabasePath 数据库路径
	DatabasePath string
	// JWTSecret JWT密钥
	JWTSecret string
	// JWTExpire JWT过期时间
	JWTExpire time.Duration
	// SystemURL 系统URL
	SystemURL string
	// CZLClientID CZL Connect 客户端ID
	CZLClientID string
	// CZLClientSecret CZL Connect 客户端密钥
	CZLClientSecret string
	// CZLAuthURL CZL Connect 授权URL
	CZLAuthURL string
	// CZLTokenURL CZL Connect 令牌URL
	CZLTokenURL string
	// CZLUserinfoURL CZL Connect 用户信息URL
	CZLUserinfoURL string
	// CZLRedirectURL CZL Connect 重定向URL
	CZLRedirectURL string
	// AIURL AI API的URL
	AIURL string
	// AIAPIKey AI API密钥
	AIAPIKey string
	// AIModel AI模型名称
	AIModel string
)

// LoadConfig 加载配置
func LoadConfig() error {
	// 加载.env文件
	godotenv.Load()

	// 应用环境
	AppEnv = getEnv("APP_ENV", "development")
	// 服务端口
	ServerPort = getEnv("SERVER_PORT", "8080")
	// 数据库路径
	DatabasePath = getEnv("DATABASE_PATH", "data/database.db")
	// JWT密钥 - 使用固定值避免重启后失效
	JWTSecret = getEnv("JWT_SECRET", "vps_monitor_secure_jwt_secret_key_2024")
	// JWT过期时间
	JWTExpire = 30 * 24 * time.Hour // 固定30天

	// 部署的域名
	SystemURL = getEnv("SYSTEM_URL", "http://localhost:3000")

	// CZL Connect 配置
	CZLClientID = getEnv("CZL_CLIENT_ID", "client_52xxx")
	CZLClientSecret = getEnv("CZL_CLIENT_SECRET", "a6d97327axxx19f9517")
	CZLAuthURL = "https://connect.czl.net/oauth2/authorize"
	CZLTokenURL = "https://connect.czl.net/api/oauth2/token"
	CZLUserinfoURL = "https://connect.czl.net/api/oauth2/userinfo"
	CZLRedirectURL = SystemURL + "/api/auth/callback"

	// AI配置
	AIURL = getEnv("AI_URL", "https://api.openai.com/v1")
	AIAPIKey = getEnv("AI_API_KEY", "")
	AIModel = getEnv("AI_MODEL", "gpt-3.5-turbo")

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
