package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(userUsecase domain.IUserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}


func (uc *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Identifier == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Identifier and password are required"})
		return
	}

	account, token, refreshToken, err := uc.userUsecase.Login(ctx, req.Identifier, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	loginResponse := &LoginResponse{
		User:        *account,
		Token:      token,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, loginResponse)
}