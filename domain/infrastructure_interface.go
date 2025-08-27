package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type IAuthMiddleware interface {
	AuthMiddleware(jwtService IJWTService) gin.HandlerFunc
	ProOnlyMiddleware() gin.HandlerFunc
	RequireRole(roles ...Role) gin.HandlerFunc
}

type IEmailService interface {
	SendPasswordResetEmail(toEmail, username, resetToken string) error
	SendVerificationEmail(toEmail, username, activationToken string) error
}

type IAIService interface {
	GenerateCompletion(ctx context.Context, prompt string) (string, error)
}

type GoogleUserInfo struct {
	ID                string
	Email             string
	Name              string
	ProfilePictureURL string
}

type IGoogleOAuthService interface {
	ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error)
}

type JWTClaims struct {
	UserID       string       `json:"user_id"`
	Role         Role         `json:"role"`
	Subscription Subscription `json:"subscription"`
	jwt.RegisteredClaims
}

type IJWTService interface {
	GenerateAccessToken(userID string, role Role) (string, *JWTClaims, error)
	GenerateRefreshToken(userID string) (string, *JWTClaims, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	ParseExpiredToken(tokenString string) (*JWTClaims, error)
	GetRefreshTokenExpiry() time.Duration
}

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type IRateLimiter interface {
	LimiterMiddleware(limit int64, period time.Duration, userIDKey string) gin.HandlerFunc
}

type ICacheService interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	AddToSet(ctx context.Context, key string, members ...any) error
	GetSetMembers(ctx context.Context, key string) ([]string, error)
	DeleteKeys(ctx context.Context, keys []string) error
}
