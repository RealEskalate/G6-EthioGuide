package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase      domain.IUserUsecase
	checklistUsecase domain.IChecklistUsecase
	refreshTokenTTL  int
}

func NewUserController(uc domain.IUserUsecase, cl domain.IChecklistUsecase, refreshTokenTTL time.Duration) *UserController {
	return &UserController{
		userUsecase:      uc,
		checklistUsecase: cl,
		refreshTokenTTL:  int(refreshTokenTTL.Seconds()),
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

func (ctrl *UserController) HandleCreateChecklist(c *gin.Context) {
	var req CreateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	user_id, err := c.Get("user_id")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are logged out, try to log in again"})
		return
	}

	userProcedure, errCreate := ctrl.checklistUsecase.CreateChecklist(c.Request.Context(), user_id.(string), req.ProcedureID)
	if errCreate != nil {
		HandleError(c, errCreate)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": ToControllerUserProcedure(userProcedure)})
}

func (ctrl *UserController) HandleGetProcedures(c *gin.Context) {
	user_id, err := c.Get("user_id")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are logged out, try to log in again"})
		return
	}

	userProcedures, errGet := ctrl.checklistUsecase.GetProcedures(c.Request.Context(), user_id.(string))
	if errGet != nil {
		HandleError(c, errGet)
		return
	}

	ProcdeureResponses := make([]*UserProcedureResponse, len(userProcedures))
	for i, prod := range userProcedures {
		ProcdeureResponses[i] = ToControllerUserProcedure(prod)
	}

	c.JSON(http.StatusOK, gin.H{"message": ProcdeureResponses})
}

func (ctrl *UserController) HandleGetChecklistById(c *gin.Context) {
	var req GetChecklistByID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	checklists, err := ctrl.checklistUsecase.GetChecklistByUserProcedureID(c.Request.Context(), req.UserProcedureID)
	if err != nil {
		HandleError(c, err)
		return
	}

	checklistResponses := make([]*ChecklistResponse, len(checklists))
	for i, check := range checklists {
		checklistResponses[i] = ToControllerChecklist(check)
	}

	c.JSON(http.StatusOK, gin.H{"message": checklistResponses})
}

func (ctrl *UserController) HandleUpdateChecklist(c *gin.Context) {
	var req UpdateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	checklist, err := ctrl.checklistUsecase.UpdateChecklist(c.Request.Context(), req.ChecklistID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": ToControllerChecklist(checklist)})
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
