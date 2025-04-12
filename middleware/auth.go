package middleware

import (
	"encoding/json"
	"fmt"
	"go-nextjs/config"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// TokenResponse OAuth2 token响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// UserInfo CZL Connect用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// Claims JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// AuthRequired 需要认证的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 检查Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		token := parts[1]
		// 验证JWT token
		claims, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("user", map[string]interface{}{
			"id":    claims.UserID,
			"email": claims.Email,
			"role":  claims.Role,
		})
		c.Next()
	}
}

// AdminRequired 需要管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 验证是否是管理员
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetLoginURL 获取登录URL
func GetLoginURL(redirectURI string) string {
	// 使用固定值，避免配置问题
	authURL := "https://connect.czl.net/oauth2/authorize"
	clientID := "client_52xxx9"

	return authURL + "?client_id=" + clientID +
		"&response_type=code" +
		"&redirect_uri=" + redirectURI +
		"&scope=openid+profile+email"
}

// GetLoginURLWithState 获取带状态的登录URL
func GetLoginURLWithState(redirectURI string, state string) string {
	// 使用固定值，避免配置问题
	authURL := "https://connect.czl.net/oauth2/authorize"
	clientID := "client_52xxx869"

	return authURL + "?client_id=" + clientID +
		"&response_type=code" +
		"&redirect_uri=" + redirectURI +
		"&scope=openid+profile+email" +
		"&state=" + state
}

// HandleCallback 处理OAuth2回调
func HandleCallback(code string, redirectURI string) (string, error) {
	// 添加调试日志
	fmt.Printf("开始处理回调，code: %s, redirectURI: %s\n", code, redirectURI)

	// 1. 使用code获取access token
	tokenResp, err := getAccessToken(code, redirectURI)
	if err != nil {
		fmt.Printf("获取access token失败: %v\n", err)
		return "", fmt.Errorf("获取access token失败: %v", err)
	}

	if tokenResp == nil || tokenResp.AccessToken == "" {
		return "", fmt.Errorf("获取到的token为空")
	}

	fmt.Printf("获取access token成功: %s...\n", tokenResp.AccessToken[:10]) // 只显示token的前几位

	// 2. 使用access token获取用户信息
	userInfo, err := getUserInfo(tokenResp.AccessToken)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		return "", fmt.Errorf("获取用户信息失败: %v", err)
	}

	if userInfo == nil {
		return "", fmt.Errorf("获取到的用户信息为空")
	}

	fmt.Printf("获取用户信息成功: %s, %s\n", userInfo.Username, userInfo.Email)

	// 3. 生成JWT token
	token, err := generateToken(userInfo)
	if err != nil {
		fmt.Printf("生成token失败: %v\n", err)
		return "", fmt.Errorf("生成token失败: %v", err)
	}
	fmt.Printf("生成JWT token成功，长度: %d\n", len(token))

	return token, nil
}

// getAccessToken 获取access token
func getAccessToken(code, redirectURI string) (*TokenResponse, error) {
	// 添加调试信息
	fmt.Printf("开始获取access token，code: %s, redirectURI: %s\n", code, redirectURI)

	// 使用固定值而不是配置，确保一致性
	url := "https://connect.czl.net/api/oauth2/token"
	clientID := "client_52xxx869"
	clientSecret := "a6d9732x"
	// 构建请求体
	rawBody := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, redirectURI, clientID, clientSecret)

	fmt.Printf("发送请求体: %s\n", rawBody)
	fmt.Printf("请求URL: %s\n", url)

	// 创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(rawBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应体: %s\n", string(respBody))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求返回非200状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 验证token
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("返回的access_token为空")
	}

	return &tokenResp, nil
}

// getUserInfo 获取用户信息
func getUserInfo(accessToken string) (*UserInfo, error) {
	// 固定的用户信息URL
	userInfoURL := "https://connect.czl.net/api/oauth2/userinfo"

	// 创建请求
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加Authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求用户信息失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("用户信息响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("用户信息响应体: %s\n", string(respBody))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("用户信息请求返回非200状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var userInfo UserInfo
	if err := json.Unmarshal(respBody, &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	// 验证必要字段
	if userInfo.ID == 0 {
		return nil, fmt.Errorf("返回的用户ID为空")
	}

	// 打印解析后的用户信息
	fmt.Printf("成功解析用户信息: ID=%d, Username=%s, Email=%s\n",
		userInfo.ID, userInfo.Username, userInfo.Email)

	return &userInfo, nil
}

// generateToken 生成JWT token
func generateToken(userInfo *UserInfo) (string, error) {
	claims := Claims{
		UserID: fmt.Sprintf("%d", userInfo.ID), // 将uint转为string
		Email:  userInfo.Email,
		Role:   "admin", // 所有认证用户都是管理员
		StandardClaims: jwt.StandardClaims{
			// 设置较长的过期时间
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(), // 30天
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// validateToken 验证JWT token
func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的token")
}
