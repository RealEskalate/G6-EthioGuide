package infrastructure

import (
	"EthioGuide/domain"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secretKey       string
	issuer          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	utilityTokenTTL time.Duration
}

// NewJWTService creates a new JWT service instance.
func NewJWTService(secret, issuer string, accessTokenTTL, refreshTokenTTL, utilityTokenTTL time.Duration) domain.IJWTService {
	return &jwtService{
		secretKey:       secret,
		issuer:          issuer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		utilityTokenTTL: utilityTokenTTL,
	}
}

func (s *jwtService) GenerateAccessToken(userID string, role domain.Role) (string, *domain.JWTClaims, error) {
	claims := &domain.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	return tokenString, claims, err
}

func (s *jwtService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// GenerateRefreshToken creates a long-lived refresh token.
func (s *jwtService) GenerateRefreshToken(userID string) (string, *domain.JWTClaims, error) {
	claims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate((time.Now().Add((s.refreshTokenTTL)))),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	return tokenString, claims, err
}

func (s *jwtService) GenerateUtilityToken(userID string) (string, *domain.JWTClaims, error) {
	claims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate((time.Now().Add((s.utilityTokenTTL)))),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	return tokenString, claims, err
}

func (s *jwtService) ParseExpiredToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if claims, ok := token.Claims.(*domain.JWTClaims); ok {
		if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
			return nil, err
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (s *jwtService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshTokenTTL
}
