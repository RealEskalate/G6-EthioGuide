package controller

import (
	"EthioGuide/domain"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase      domain.IUserUsecase
	procedureUsecase domain.ISearchUseCase
	refreshTokenTTL  int
}

func NewUserController(uc domain.IUserUsecase, puc domain.ISearchUseCase, refreshTokenTTL time.Duration) *UserController {
	return &UserController{
		userUsecase:      uc,
		procedureUsecase: puc,
		refreshTokenTTL:  int(refreshTokenTTL.Seconds()),
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

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": ToUserResponse(account)})
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
			User:         ToUserResponse(account),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		setAuthCookie(c, refreshToken, ctrl.refreshTokenTTL)
		c.JSON(http.StatusOK, &LoginResponse{
			User:        ToUserResponse(account),
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

	c.JSON(http.StatusOK, ToUserResponse(account))
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
			User:         ToUserResponse(account),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		setAuthCookie(c, refreshToken, ctrl.refreshTokenTTL)
		c.JSON(http.StatusOK, &LoginResponse{
			User:        ToUserResponse(account),
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

	c.JSON(http.StatusOK, ToUserResponse(savedAccount))
}

// @Summary      Logout
// @Description  Logout a user.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200 {string}  "Log out successful"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /auth/logout/ [post]
func (ctrl *UserController) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	err := ctrl.userUsecase.Logout(c.Request.Context(), userID.(string))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// @Summary      Forgot Password
// @Description  Forgot password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ForgotDTO true "User Email"
// @Success      200 {string}  "Reset token sent"
// @Failure      400 {string}  "Invalid request"
// @Failure      500 {string}  "Server error"
// @Router       /auth/forgot [post]
func (ctrl *UserController) HandleForgot(c *gin.Context) {
	var req ForgotDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := ctrl.userUsecase.ForgetPassword(c.Request.Context(), req.Email)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset token sent"})
}

// @Summary      Reset Password
// @Description  Reset password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ResetDTO true "Reset Token and New Password"
// @Success      200 {string}  "Reset token sent"
// @Failure      400 {string}  "Invalid request"
// @Failure      500 {string}  "Server error"
// @Router       /auth/reset [post]
func (ctrl *UserController) HandleReset(c *gin.Context) {
	var req ResetDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := ctrl.userUsecase.ResetPassword(c.Request.Context(), req.ResetToken, req.NewPassword)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password Updated Successfully"})
}

// @Summary      Verify Account
// @Description  Verify User Account.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ActivateDTO true "Reset Token and New Password"
// @Success      200 {string}  "Reset token sent"
// @Failure      400 {string}  "Invalid request"
// @Failure      500 {string}  "Server error"
// @Router       /auth/verify [post]
func (ctrl *UserController) HandleVerify(c *gin.Context) {
	var req ActivateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := ctrl.userUsecase.VerifyAccount(c.Request.Context(), req.ActivateToken); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Activated Successfully"})
}

func (ctrl *UserController) HandleCreateOrg(c *gin.Context) {
	var req OrgCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := ctrl.userUsecase.RegisterOrg(c.Request.Context(), req.Name, req.Email, req.OrgType)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organization created Successsfully"})
}

func (ctrl *UserController) HandleGetOrgs(c *gin.Context) {
	var filter domain.GetOrgsFilter
	filter.Type = c.Query("type")
	filter.Query = c.Query("q")

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "1"), 10, 64)

	filter.Page = page
	filter.PageSize = pageSize

	accounts, total, err := ctrl.userUsecase.GetOrgs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not get organizations"})
		return
	}

	orgsResponse := make([]OrganizationResponseDTO, len(accounts))
	for i, acc := range accounts {
		orgsResponse[i] = ToOrganizationDTO(acc)
	}

	pagination := &PaginatedOrgsResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	c.JSON(http.StatusOK, gin.H{"data": orgsResponse, "pagination": pagination})
}

func (ctrl *UserController) HandleGetOrgById(c *gin.Context) {
	orgId := c.Param("id")
	account, err := ctrl.userUsecase.GetOrgById(c.Request.Context(), orgId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ToOrganizationDetailDTO(account)})
}

func (ctrl *UserController) HandleUpdateOrgs(c *gin.Context) {
	orgId := c.Param("id")
	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	update := make(map[string]interface{})

	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.ProfilePicURL != nil {
		update["profile_pic_url"] = *req.ProfilePicURL
	}
	if req.Description != nil {
		update["organization_detail.description"] = *req.Description
	}
	if req.Location != nil {
		update["organization_detail.location"] = *req.Location
	}
	if req.PhoneNumbers != nil {
		update["organization_detail.phone_numbers"] = req.PhoneNumbers
	}
	if req.ContactInfo != nil {
		if req.ContactInfo.Website != nil {
			update["organization_detail.contact_info.website"] = *req.ContactInfo.Website
		}
		if req.ContactInfo.Socials != nil {
			update["organization_detail.contact_info.socials"] = req.ContactInfo.Socials
		}
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	err := ctrl.userUsecase.UpdateOrgFields(c.Request.Context(), orgId, update)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organization updated successfully"})
}

func (ctrl *UserController) HandleSearch(c *gin.Context) {
	query := c.Param("q")
	page, err := strconv.ParseInt(c.Param("page"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	limit, errlimit := strconv.ParseInt(c.Param("limit"), 10, 64)
	if errlimit != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result := &domain.SearchFilterRequest{
		Query: query,
		Page:  page,
		Limit: limit,
	}

	searchResult, err := ctrl.procedureUsecase.Search(c.Request.Context(), *result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to search"})
		return
	}

	c.JSON(http.StatusOK, ToSearchJSON(searchResult))
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
