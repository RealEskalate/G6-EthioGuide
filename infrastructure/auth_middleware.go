package infrastructure

import (
	"EthioGuide/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a gin middleware for JWT authentication.
func AuthMiddleware(jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized,  gin.H{"error":"Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized,gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{ "error": "invalid or expired token"})
			return
		}

		// Set user info into the context for later use in handlers
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("userSubscription", claims.Subscription)

		c.Next()
	}
}

func ProOnlyMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		userSubscription, exists := c.Get("userSubscription")
		if !exists{
			c.JSON(http.StatusForbidden, gin.H{"error": "Subscription information is missing"})
			return
		}

		subscription, ok := userSubscription.(string)
		if !ok {
			c.JSON(http.StatusForbidden,gin.H{"error": "Invalid subscription type"})
			return
		}

		if subscription == "pro"{
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden,gin.H{"error": "Insufficient subscription permissions"})

	}
	
}


// RequireRole restricts access to specified roles.
func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, "Role information is missing")
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role type"})
			return
		}

		for _, allowed := range roles {
			if domain.Role(role) == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role permissions"})
	}
}
