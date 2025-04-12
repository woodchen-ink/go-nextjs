package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username  string `gorm:"size:100;not null;unique"` // 用户名
	Nickname  string `gorm:"size:100"`                 // 昵称
	Email     string `gorm:"size:100"`                 // 邮箱
	Avatar    string `gorm:"size:500"`                 // 头像
	Role      string `gorm:"size:20;default:user"`     // 角色：admin, user
	ExternID  string `gorm:"size:100;unique"`          // 外部ID
	Provider  string `gorm:"size:20"`                  // 提供商：czl_connect
	LastLogin int64  `gorm:"default:0"`                // 最后登录时间
	Token     string `gorm:"size:500"`                 // 访问令牌
}

// UserInfo 用户信息响应
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}
