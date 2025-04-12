package handler

import (
	"fmt"
	"go-nextjs/config"
	"go-nextjs/middleware"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Login 处理登录请求
func Login(c *gin.Context) {
	state := c.Query("state") // 获取state参数，用于保存回调后的目标URL

	// 使用固定的后端回调URL
	redirectURI := config.CZLRedirectURL

	// 生成授权URL
	var loginURL string
	if state != "" {
		loginURL = middleware.GetLoginURLWithState(redirectURI, state)
	} else {
		loginURL = middleware.GetLoginURL(redirectURI)
	}

	// 记录生成的URL，用于调试
	fmt.Printf("生成的授权URL: %s\n", loginURL)

	// 返回授权URL
	c.JSON(http.StatusOK, gin.H{"url": loginURL})
}

// Callback 处理OAuth2回调
func Callback(c *gin.Context) {
	// 记录回调参数
	fullURL := c.Request.URL.String()
	fmt.Printf("收到回调请求，完整URL: %s\n", fullURL)

	code := c.Query("code")
	if code == "" {
		fmt.Println("回调错误: 未提供授权码")
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供授权码"})
		return
	}
	fmt.Printf("收到授权码，长度: %d\n", len(code))

	// 获取state参数，这是用于防止CSRF攻击的随机值
	state := c.Query("state")
	fmt.Printf("回调state参数: %s\n", state)

	// 使用与授权请求相同的重定向URL
	redirectURI := config.CZLRedirectURL
	fmt.Printf("使用重定向URL: %s\n", redirectURI)

	// 处理OAuth2回调
	token, err := middleware.HandleCallback(code, redirectURI)
	if err != nil {
		fmt.Printf("处理回调失败: %v\n", err)
		// 前端回调页面
		frontendRedirectURL := config.SystemURL + "/auth/callback"
		// 重定向回前端，携带错误信息和原始state
		errorRedirectURL := frontendRedirectURL + "?error=" + url.QueryEscape("处理回调失败: "+err.Error())
		if state != "" {
			errorRedirectURL += "&state=" + url.QueryEscape(state)
		}
		c.Redirect(http.StatusTemporaryRedirect, errorRedirectURL)
		return
	}

	// 验证生成的token不为空
	if token == "" {
		fmt.Println("回调错误: 生成的token为空")
		frontendRedirectURL := config.SystemURL + "/auth/callback"
		errorRedirectURL := frontendRedirectURL + "?error=" + url.QueryEscape("生成的token为空")
		if state != "" {
			errorRedirectURL += "&state=" + url.QueryEscape(state)
		}
		c.Redirect(http.StatusTemporaryRedirect, errorRedirectURL)
		return
	}

	fmt.Printf("成功生成token，长度: %d\n", len(token))

	// 前端回调页面URL
	frontendRedirectURL := config.SystemURL + "/auth/callback"
	// 重定向到前端，并携带token和state
	redirectURL := frontendRedirectURL + "?token=" + url.QueryEscape(token)
	if state != "" {
		redirectURL += "&state=" + url.QueryEscape(state)
	}

	fmt.Printf("最终重定向URL: %s\n", redirectURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout 处理登出请求
func Logout(c *gin.Context) {
	// 由于使用的是JWT，服务端不需要特殊处理，只需前端清除token
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
