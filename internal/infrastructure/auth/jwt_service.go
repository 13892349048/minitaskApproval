package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/pkg/errors"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// JWTService JWT服务实现
type JWTService struct {
	config shared.JWTConfig
}

// NewJWTService 创建JWT服务
func NewJWTService(config shared.JWTConfig) shared.JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func (j *JWTService) GenerateTokens(userID, email string, roles []string) (*shared.TokenPair, error) {
	now := time.Now()

	// 生成访问令牌
	accessToken, err := j.generateToken(userID, email, roles, shared.TokenTypeAccess, now.Add(j.config.AccessTokenExpiry))
	if err != nil {
		logger.Error("Failed to generate access token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := j.generateToken(userID, email, roles, shared.TokenTypeRefresh, now.Add(j.config.RefreshTokenExpiry))
	if err != nil {
		logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &shared.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(j.config.AccessTokenExpiry.Seconds()),
		ExpiresAt:    now.Add(j.config.AccessTokenExpiry),
	}, nil
}

// ValidateToken 验证访问令牌
func (j *JWTService) ValidateToken(tokenString string) (*shared.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shared.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		// JWT v5版本的错误处理
		if err.Error() == "token is expired" {
			return nil, errors.ErrExpiredToken
		}
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*shared.Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	// 验证令牌类型
	if claims.TokenType != shared.TokenTypeAccess {
		return nil, errors.ErrInvalidTokenType
	}

	return claims, nil
}

// RefreshToken 使用刷新令牌生成新的访问令牌
func (j *JWTService) RefreshToken(refreshToken string) (*shared.TokenPair, error) {
	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(refreshToken, &shared.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*shared.Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	// 验证是否是刷新令牌
	if claims.TokenType != shared.TokenTypeRefresh {
		return nil, errors.ErrInvalidTokenType
	}

	// 生成新的令牌对
	return j.GenerateTokens(claims.UserID, claims.Email, claims.Roles)
}

// RevokeToken 撤销令牌（简单实现，生产环境可考虑使用Redis黑名单）
func (j *JWTService) RevokeToken(tokenString string) error {
	// 简单实现：记录日志
	// 生产环境可以维护一个黑名单（使用Redis）
	logger.Info("Token revoked", zap.String("token", tokenString[:20]+"..."))

	// TODO: 实现令牌黑名单
	// 可以使用Redis存储被撤销的令牌，直到它们过期

	return nil
}

// generateToken 生成JWT令牌
func (j *JWTService) generateToken(userID, email string, roles []string, tokenType string, expiresAt time.Time) (string, error) {
	claims := shared.Claims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// 为什么这样实现？
//
// 1. 安全性考虑：
//    - 使用HMAC-SHA256签名算法（安全且高效）
//    - 验证签名方法，防止算法替换攻击
//    - 区分访问令牌和刷新令牌类型
//
// 2. 错误处理：
//    - 区分不同类型的验证错误（过期、无效、类型错误）
//    - 详细的错误日志记录
//    - 返回标准化的错误类型
//
// 3. 令牌刷新机制：
//    - 刷新令牌验证后生成新的令牌对
//    - 保持用户信息的一致性
//    - 自动延长会话时间
//
// 4. 可扩展性：
//    - RevokeToken为黑名单机制预留接口
//    - Claims结构可以轻松扩展更多用户信息
//    - 配置化的过期时间和签名密钥
