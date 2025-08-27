package domain

import "time"

type TokenType string

const (
	AccessToken        TokenType = "access token"
	RefreshToken       TokenType = "refresh token"
	VerificationToken  TokenType = "verification token"
	ResetPasswordToken TokenType = "reset password token"
)

type Token struct {
	Id        string
	Token     string
	TokenType TokenType
	ExpiresAt time.Time
}
