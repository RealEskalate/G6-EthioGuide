package router

import (
	"EthioGuide/delivery/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and registers all application routes.
func SetupRouter(
	userController *controller.UserController,
	geminiController *controller.GeminiController,
	aiChatController *controller.AIChatController,
	authMiddleware gin.HandlerFunc,
	proOnlyMiddleware gin.HandlerFunc,
	requireAdminRole gin.HandlerFunc,
) *gin.Engine {

	router := gin.Default()

	// Health check endpoint - always public
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Group all routes under a versioned prefix
	v1 := router.Group("/api/v1")
	{
		// --- Public Routes ---
		// These endpoints do not require any authentication.
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userController.Register)
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/refresh-token", userController.HandleRefreshToken)
		}

		// --- Private Routes (Require Authentication) ---
		// All routes in this group are protected by the base AuthMiddleware.
		apiGroup := v1.Group("/")
		apiGroup.Use(authMiddleware)
		{
			// --- Standard User Routes ---
			// Any logged-in user (regardless of role or subscription) can access these.
			aiGroup := apiGroup.Group("/ai")
			aiGroup.Use(authMiddleware)
			{
				aiGroup.POST("/translate", geminiController.Translate)
				aiGroup.POST("/guide", aiChatController.AIChatController)
			}

			// --- PRO Subscription Routes ---
			// These routes require the user to be logged in AND have a "pro" subscription.
			// We chain the ProOnlyMiddleware after the main auth middleware.
			proGroup := apiGroup.Group("/pro")
			proGroup.Use(proOnlyMiddleware)
			{
				// Example: A more advanced AI feature only for paying users.
				// You would need to create this controller method.
				// proGroup.POST("/ai/summarize", geminiController.SummarizeContent)
			}

			// --- Admin Routes ---
			// These routes require the user to be logged in AND have the "Admin" role.
			adminGroup := apiGroup.Group("/admin")
			adminGroup.Use(requireAdminRole)
			{
				// Example: An endpoint for an admin to get a list of all users.
				// You would need to create this controller method.
				// adminGroup.GET("/users", userController.GetAllUsers)
			}
		}
	}

	return router
}
