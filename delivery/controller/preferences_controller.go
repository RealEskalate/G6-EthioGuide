package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PreferencesController struct {
	prerfUsecase domain.IPreferencesUsecase
}

func NewPreferencesController(prefUsecase domain.IPreferencesUsecase) *PreferencesController{
	return &PreferencesController{
		prerfUsecase: prefUsecase,
	}
}

// @Summary      Get Preferences
// @Description  Get user preferences
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  PreferencesDTO "Detail of user preferences"
// @Failure      400  {string}   "Bad Request"
// @Failure      500  {string}   "Server error"
// @Router       /auth/me/preferences [get]
func (P *PreferencesController) GetUserPreferences(c *gin.Context) {
    // Assume userID is extracted from context/session/JWT
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    pref, err := P.prerfUsecase.GetUserPreferences(c.Request.Context(), userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Preferences not found"})
        return
    }
    dto := PreferencesDTO{
        PreferredLang:     string(pref.PreferredLang),
        PushNotification:  pref.PushNotification,
        EmailNotification: pref.EmailNotification,
    }
    c.JSON(http.StatusOK, dto)
}

// @Summary      Update Preferences
// @Description  Update user preferences
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body PreferencesDTO true "update user preferences detail."
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  PreferencesDTO "Detail of user preferences"
// @Failure      400  {string}   "Bad Request"
// @Failure      500  {string}   "Server error"
// @Router       /auth/me/preferences [patch]
func (P *PreferencesController) UpdateUserPreferences(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    var dto PreferencesDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    pref := &domain.Preferences{
        UserID:            userID.(string),
        PreferredLang:     domain.Lang(dto.PreferredLang),
        PushNotification:  dto.PushNotification,
        EmailNotification: dto.EmailNotification,
    }
    if err := P.prerfUsecase.UpdateUserPreferences(c.Request.Context(), pref); err != nil {
        HandleError(c, err)
        return
    }
    c.JSON(http.StatusOK, dto)
}
