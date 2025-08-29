package controller

import (
	"EthioGuide/domain"
	"errors"
	"net/http"
	"strconv"
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
