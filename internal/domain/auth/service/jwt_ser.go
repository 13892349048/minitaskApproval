package service

import "github.com/taskflow/internal/domain/auth/valueobject"

// JWTService JWT服务接口
type JWTService interface {
	// GenerateTokens 生成访问令牌和刷新令牌
	GenerateTokens(userID, email string, roles []string) (*valueobject.TokenPair, error)

	// ValidateToken 验证访问令牌
	ValidateToken(tokenString string) (*valueobject.Claims, error)

	// RefreshToken 使用刷新令牌生成新的访问令牌
	RefreshToken(refreshToken string) (*valueobject.TokenPair, error)

	// RevokeToken 撤销令牌（可选实现，用于登出）
	RevokeToken(tokenString string) error
}
