package shared

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService JWT服务接口
type JWTService interface {
	// GenerateTokens 生成访问令牌和刷新令牌
	GenerateTokens(userID, email string, roles []string) (*TokenPair, error)

	// ValidateToken 验证访问令牌
	ValidateToken(tokenString string) (*Claims, error)

	// RefreshToken 使用刷新令牌生成新的访问令牌
	RefreshToken(refreshToken string) (*TokenPair, error)

	// RevokeToken 撤销令牌（可选实现，用于登出）
	RevokeToken(tokenString string) error
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // 访问令牌过期时间（秒）
	ExpiresAt    time.Time `json:"expires_at"` // 访问令牌过期时间点
}

// Claims JWT声明
type Claims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `json:"secret"`
	AccessTokenExpiry  time.Duration `json:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
	Issuer             string        `json:"issuer"`
}

// TokenType 令牌类型常量
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// 为什么这样设计？
//
// 1. 双令牌机制：
//    - Access Token：短期有效（24小时），用于API访问
//    - Refresh Token：长期有效（7天），用于刷新Access Token
//    - 安全性：即使Access Token被盗，影响时间有限
//
// 2. Claims结构设计：
//    - 包含必要的用户信息：ID、邮箱、角色
//    - 使用标准的RegisteredClaims
//    - TokenType区分访问令牌和刷新令牌
//
// 3. 接口设计：
//    - GenerateTokens：登录时生成令牌对
//    - ValidateToken：中间件验证令牌
//    - RefreshToken：自动刷新机制
//    - RevokeToken：登出时撤销令牌（可选）
