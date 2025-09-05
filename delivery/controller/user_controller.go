package controller

import (
	"EthioGuide/domain"
	"net/http"
	"net/url"
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

// @Summary      Refresh Access Token
// @Description  Refresh Access Token for web and mobile
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200 {string}  "New Access Token"
// @Failure      400 {string}  "invalid
// @Failure      409 {string}  "invalid"
// @Failure      500 {string}  "invalid"
// @Router       /auth/refresh [post]
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

// @Summary      Register a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User Registration Details"
// @Success      201 {object} UserResponse "User created Successfully"
// @Failure      400 {string}  "invalid
// @Failure      409 {string}  "invalid"
// @Failure      500 {string}  "invalid"
// @Router       /auth/register [post]
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

// @Summary      Login a new user
// @Description  Login a user account with the provided details.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User Registration Details"
// @Success      201 {object} LoginResponse  "Login Successful"
// @Failure      400 {string}  "invalid
// @Failure      500 {string}  "invalid
// @Router       /auth/login [post]
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
			User:         toUserResponse(account),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		setAuthCookie(c, refreshToken, ctrl.refreshTokenTTL)
		c.JSON(http.StatusOK, &LoginResponse{
			User:        toUserResponse(account),
			AccessToken: accessToken,
		})
	}
}

// @Summary      Get Profile
// @Description  Get user's profile detail.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200 {object} UserResponse "Profile Retrieved"
// @Failure      400 {string}  "Invalid request"
// @Failure      404 {string}  "Invalid request"
// @Failure      500 {string}  "Server error"
// @Router       /auth/me [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	account, err := ctrl.userUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toUserResponse(account))
}

// @Summary      Update password
// @Description  Update user's password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body ChangePasswordRequest true "Password Change Detail"
// @Success      200 {string}  "Password changed"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /auth/me/password [patch]
func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	accountID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := ctrl.userUsecase.UpdatePassword(c.Request.Context(), accountID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// @Summary      Social Login
// @Description  Login with third party auth.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body SocialLoginRequest true "Social Login Detail."
// @Success      200 {object} LoginResponse "Login successful"
// @Failure      400 {string}  "Invalid request"
// @Failure      500 {string}  "Server error"
// @Router       /auth/social [post]
func (ctrl *UserController) SocialLogin(c *gin.Context) {
	var req SocialLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	code, err := url.QueryUnescape(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to URL-decode the pasted code: " + err.Error()})
		return
	}

	account, accessToken, refreshToken, err := ctrl.userUsecase.LoginWithSocial(c.Request.Context(), req.Provider, code)
	if err != nil {
		HandleError(c, err)
		return
	}

	if isMobileClient(c) {
		c.JSON(http.StatusOK, &LoginResponse{
			User:         toUserResponse(account),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		setAuthCookie(c, refreshToken, ctrl.refreshTokenTTL)
		c.JSON(http.StatusOK, &LoginResponse{
			User:        toUserResponse(account),
			AccessToken: accessToken,
		})
	}
}

// @Summary      Update Profile
// @Description  Update user's profile.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body UserUpdateRequest true "Updated Account Details"
// @Success      200 {object} domain.Account "Account Updated"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /auth/me/ [patch]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
        return
    }

    var req UserUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 1. Fetch current account from usecase
    account, err := ctrl.userUsecase.GetProfile(c.Request.Context(), userID.(string))
    if err != nil {
        HandleError(c, err)
        return
    }

    // 2. Convert DTO → domain.Account with updates applied
    updatedAccount := ToDomainAccountUpdate(&req, account)

    // 3. Call usecase with pure domain model
    savedAccount, err := ctrl.userUsecase.UpdateProfile(c.Request.Context(), updatedAccount)
    if err != nil {
        HandleError(c, err)
        return
    }

    c.JSON(http.StatusOK, toUserResponse(savedAccount))
}

func (ctrl *UserController) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	err := ctrl.userUsecase.Logout(c.Request.Context(), userID.(string))
	if err != nil{
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
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
