package controller

import (
	"EthioGuide/domain"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	log.Panicf("Error %v\n", err)
	switch {
	// --- 400 Bad Request ---
	// Catch specific validation errors first
	case errors.Is(err, domain.ErrPasswordTooShort),
		errors.Is(err, domain.ErrInvalidEmailFormat),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrInvalidProvider),
		errors.Is(err, domain.ErrUsernameEmpty),
		errors.Is(err, domain.ErrUsernameTooLong),
		errors.Is(err, domain.ErrInvalidBody),
		errors.Is(err, domain.ErrUnsupportedLanguage):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	// Catch generic validation error
	case errors.Is(err, domain.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input provided"})

	// --- 401 Unauthorized ---
	case errors.Is(err, domain.ErrAuthenticationFailed),
		errors.Is(err, domain.ErrInvalidActivationToken),
		errors.Is(err, domain.ErrInvalidResetToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	// --- 403 Forbidden ---
	case errors.Is(err, domain.ErrPermissionDenied),
		errors.Is(err, domain.ErrCannotChangeOwnRole), // Specific forbidden action
		errors.Is(err, domain.ErrOAuthUser),           // Specific forbidden action
		errors.Is(err, domain.ErrAccountNotActive):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	// --- 404 Not Found ---
	case errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	// Add other "not found" errors here if you have them (e.g., domain.ErrBlogNotFound)

	// --- 409 Conflict ---
	case errors.Is(err, domain.ErrEmailExists),
		errors.Is(err, domain.ErrUsernameExists),
		errors.Is(err, domain.ErrPhoneNumberExists),
		errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	// --- 500 Internal Server Error (Default) ---
	default:
		log.Printf("Internal Server Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected internal error occurred. Please try again later."})
	}
}
