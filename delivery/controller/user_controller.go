package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strings"

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

func setAuthCookie(c *gin.Context, accessToken string) {
	if accessToken != "" {
		c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)
	}
}

func (ac *AuthController) HandleRefreshToken(c *gin.Context) {
	userAgent := c.Request.UserAgent()
	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "Android") {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		refreshToken := strings.TrimPrefix(authHeader, prefix)
		newAccess, newRefresh, err := ac.authusecase.RefreshTokenForMobile(c.Request.Context(), refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		})
		
	} else {
		refreshtoken, err := c.Cookie("refresh_token")
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		}

		newAccess, err := ac.authusecase.RefreshTokenForWeb(c.Request.Context(), refreshtoken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		}

		setAuthCookie(c, newAccess)
		c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
	}

}
