package domain

import "time"

type TokenModel struct {
	Id        string
	Token     string
	TokenType string
	ExpiresAt time.Time
}
