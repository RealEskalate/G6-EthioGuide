package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authusecase domain.IAuthInterface
}

func NewAuthController(usecase domain.IAuthInterface) *AuthController {
	return &AuthController{
		authusecase: usecase,
	}
}

func setAuthCookie(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 30*24*60*60, "/", "", false, true)
}

func (ac *AuthController) HandleRefreshToken(c *gin.Context) {
	refreshtoken, err := c.Cookie("refresh_token")
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
	}

	newAccess, newRefresh, err := ac.authusecase.RefreshToken(c.Request.Context(), refreshtoken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
	}

	setAuthCookie(c, newAccess, newRefresh)
	c.JSON(http.StatusOK, gin.H{"message": "tokens refreshed"})
}
