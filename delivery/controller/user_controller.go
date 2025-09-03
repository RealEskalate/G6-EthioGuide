package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase     domain.IUserUsecase
	refreshTokenTTL int
}

func NewUserController(uc domain.IUserUsecase, refreshTokenTTL time.Duration) *UserController {
	return &UserController{
		userUsecase:     uc,
		refreshTokenTTL: int(refreshTokenTTL.Seconds()),
	}
}

func isMobileClient(c *gin.Context) bool {
	return c.GetHeader("X-Client-Type") == "mobile"
}

func (ctrl *UserController) HandleRefreshToken(c *gin.Context) {
	if isMobileClient(c) {
		// --- Mobile Client Logic ---
		refreshToken, err := extractBearerToken(c)
		if err != nil {
			HandleError(c, err)
			return
		}

		newAccess, newRefresh, err := ctrl.userUsecase.RefreshTokenForMobile(c.Request.Context(), refreshToken)
		if err != nil {
			HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		})

	} else {
		// --- Web Client Logic ---
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			HandleError(c, domain.ErrAuthenticationFailed)
			return
		}

		newAccess, err := ctrl.userUsecase.RefreshTokenForWeb(c.Request.Context(), refreshToken)
		if err != nil {
			HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
	}
}

// Register is now much cleaner, delegating all error handling to the helper.
func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userDetail := &domain.UserDetail{
		Username:         req.Username,
		SubscriptionPlan: domain.SubscriptionNone,
		IsBanned:         false,
		IsVerified:       false,
	}
	account := &domain.Account{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         domain.RoleUser,
		UserDetail:   userDetail,
	}

	err := ctrl.userUsecase.Register(c.Request.Context(), account)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": toUserResponse(account)})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	account, accessToken, refreshToken, err := ctrl.userUsecase.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	if isMobileClient(c) {
		c.JSON(http.StatusOK, &LoginResponse{
			User:         *account,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		setAuthCookie(c, refreshToken, ctrl.refreshTokenTTL)
		c.JSON(http.StatusOK, &LoginResponse{
			User:        *account,
			AccessToken: accessToken,
		})
	}
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	account, err := ctrl.userUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		HandleError(c, err)
	}

	c.JSON(http.StatusOK, toUserResponse(account))
}

func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	accountID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
	}

	err := ctrl.userUsecase.UpdatePassword(c.Request.Context(), accountID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// --- HELPER FUNCTIONS ---

// extractBearerToken is a helper to get the token from the Authorization header.
func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", domain.ErrAuthenticationFailed
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", domain.ErrAuthenticationFailed
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}

func setAuthCookie(c *gin.Context, refreshToken string, refreshTokenTTL int) {
	if refreshToken != "" {
		c.SetCookie("refresh_token", refreshToken, refreshTokenTTL, "/", "", false, true)
	}
}
