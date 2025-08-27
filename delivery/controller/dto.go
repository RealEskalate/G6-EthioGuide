package controller

import "EthioGuide/domain"

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User      domain.Account `json:"user"`
	Token     string         `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}
