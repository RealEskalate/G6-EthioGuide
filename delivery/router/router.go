package router

import (
	"EthioGuide/delivery/controller"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "EthioGuide/docs"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin router and registers all application routes.
func SetupRouter(
	userController *controller.UserController,
	procedureController *controller.ProcedureController,
	catagorieController *controller.CategoryController,
	geminiController *controller.GeminiController,
	feedbackController *controller.FeedbackController,
	postController *controller.PostController,
	authMiddleware gin.HandlerFunc,
	proOnlyMiddleware gin.HandlerFunc,
	requireAdminRole gin.HandlerFunc,
	requireAdminOrOrgRole gin.HandlerFunc,
) *gin.Engine {

	router := gin.Default()

	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://ethio-guide.vercel.app",
			"https://your-production-site.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Client-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// Health check endpoint - always public
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group all routes under a versioned prefix
	v1 := router.Group("/api/v1")
	{
		// --- Public Routes ---
		// These endpoints do not require any authentication.
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userController.Register)
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/refresh", userController.HandleRefreshToken)
			authGroup.POST("/social", userController.SocialLogin)
		}

		// --- Private Routes (Require Authentication) ---
		// All routes in this group are protected by the base AuthMiddleware.
		apiGroup := v1.Group("/")
		apiGroup.Use(authMiddleware)
		{
			// --- Standard User Routes ---
			// Any logged-in user (regardless of role or subscription) can access these.
			aiGroup := apiGroup.Group("/ai")
			{
				aiGroup.POST("/translate", geminiController.Translate)
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

			// --- User Profile Routes ---
			// Any logged in user can access these routes to manage their profile
			authGroup := apiGroup.Group("/auth")
			authGroup.Use(authMiddleware)
			{
				authGroup.GET("/me", userController.GetProfile)
				authGroup.PATCH("/me/password", userController.UpdatePassword)
				authGroup.PATCH("/me", userController.UpdateProfile)
			}

			procedures := v1.Group("/procedures")
			{
				procedures.POST("", authMiddleware, requireAdminOrOrgRole, procedureController.CreateProcedure)
				procedures.POST("/:id/feedback", authMiddleware, feedbackController.SubmitFeedback)
				procedures.GET("/:id/feedback", feedbackController.GetAllFeedbacksForProcedure)
			}

			feedback := v1.Group("/feedback")
			{
				feedback.PATCH("/:id", authMiddleware, requireAdminOrOrgRole, feedbackController.UpdateFeedbackStatus)
			}

			discussions := v1.Group("/discussions")
			{
				discussions.POST("", authMiddleware, postController.CreatePost)
				discussions.GET("", postController.GetPosts)
				discussions.GET("/:id", postController.GetPostByID)
				discussions.PATCH("/:id", authMiddleware, postController.UpdatePost)
				discussions.DELETE("/:id", authMiddleware, requireAdminOrOrgRole, postController.DeletePost)
			}

			categories := v1.Group("/categories")
			{
				categories.POST("", authMiddleware, requireAdminOrOrgRole, catagorieController.CreateCategory)
				categories.GET("", catagorieController.GetCategory)
			}

		}
	}

	// MOCK ROUTES
	{
		// 1) Auth & Accounts
		auth := v1.Group("/auth")
		{
			auth.POST("/verify", handleVerifyEmail)
			auth.POST("/forgot", handleForgot)
			auth.POST("/reset", handleReset)
		}

		// 2) Users & Profiles
		users := v1.Group("/users")
		{
			users.GET("/:id", handleGetUser)
			users.GET("/me/preferences", handleGetUserPreferences)
			users.PATCH("/me/preferences", handleUpdateUserPreferences)
			users.GET("/me/summary", handleGetUserSummary)
		}

		// 3) Organizations
		orgs := v1.Group("/orgs")
		{
			orgs.POST("", handleCreateOrg)
			orgs.GET("/pending", handleGetPendingOrgs)
			orgs.PATCH("/:id/approve", handleApproveOrg)
			orgs.GET("", handleGetOrgs)
			orgs.GET("/:id", handleGetOrg)
			orgs.PATCH("/:id", handleUpdateOrg)
			orgs.GET("/:id/feedback", handleGetOrgFeedback)
		}

		// 4) Categories & Taxonomy
		categories := v1.Group("/categories")
		{
			categories.PATCH("/:id", handleUpdateCategory)
		}

		// 5) Procedures (core)
		procedures := v1.Group("/procedures")
		{
			procedures.GET("", handleGetProcedures)
			procedures.GET("/:id", handleGetProcedure)
			procedures.PATCH("/:id", handleUpdateProcedure)
			procedures.DELETE("/:id", handleDeleteProcedure)
			procedures.PATCH("/:id/verify", handleVerifyProcedure)
			procedures.GET("/:id/audit", handleGetProcedureAudit)
			procedures.GET("/popular", handleGetPopularProcedures)
			procedures.GET("/recent", handleGetRecentProcedures)
		}

		// 6) Search & Discovery
		v1.GET("/search", handleSearch)

		// 7) Checklists & Progress
		checklists := v1.Group("/checklists")
		{
			checklists.POST("", handleCreateChecklist)
			checklists.GET("/:userProcedureId", handleGetChecklistByUserProcedureId) // Assuming this mapping
			checklists.PATCH("/:checklistID", handleUpdateChecklistItem)
		}
		v1.GET("/myProcedures", handleGetMyProcedures)

		// 8) Documents & File Vault
		v1.POST("/uploads/signature", handleUploadSignature)
		documents := v1.Group("/documents")
		{
			documents.POST("", handleCreateDocument)
			documents.GET("", handleGetDocuments)
			documents.GET("/:id", handleGetDocument)
			documents.PATCH("/:id", handleUpdateDocument)
			documents.DELETE("/:id", handleDeleteDocument)
		}

		// 9) Reminders & Notifications
		reminders := v1.Group("/reminders")
		{
			reminders.POST("", handleCreateReminder)
			reminders.GET("", handleGetReminders)
			reminders.PATCH("/:id", handleUpdateReminder)
			reminders.DELETE("/:id", handleDeleteReminder)
		}
		notifications := v1.Group("/notifications")
		{
			notifications.GET("", handleGetNotifications)
			notifications.PATCH("/:id/read", handleMarkNotificationRead)
		}

		// 10) Discussions (Community)
		discussions := v1.Group("/discussions")
		{
			discussions.POST("/:id/upvote", handleUpvoteDiscussion)
			discussions.POST("/:id/downvote", handleDownvoteDiscussion)
			discussions.POST("/:id/report", handleReportDiscussion)
		}

		// 11) Feedback (standalone updates)
		feedback := v1.Group("/feedback")
		{
			feedback.POST("/:id/upvote", handleUpvoteFeedback)
		}

		// 12) Official Notices
		notices := v1.Group("/notices")
		{
			notices.POST("", handleCreateNotice)
			notices.GET("", handleGetNotices)
			notices.GET("/:id", handleGetNotice)
			notices.PATCH("/:id", handleUpdateNotice)
			notices.DELETE("/:id", handleDeleteNotice)
		}

		// 13) AI Guidance (Gemini)
		ai := v1.Group("/ai")
		{
			ai.POST("/guide", handleAIGuide)
			ai.GET("/history", handleAIGetHistory)
			ai.POST("/mark-not-verified", handleAIMarkNotVerified)
			ai.POST("/speech-to-text", handleAISpeechToText)
		}

		// 14) Direct Messages
		dm := v1.Group("/dm/threads")
		{
			dm.POST("", handleCreateDMThread)
			dm.GET("", handleGetDMThreads)
			dm.GET("/:id", handleGetDMThread)
			dm.POST("/:id/messages", handleCreateDMMessage)
			dm.PATCH("/:id/close", handleCloseDMThread)
		}

		// 15) Subscriptions & Payments
		v1.GET("/plans", handleGetPlans)
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", handleCreateSubscription)
			subscriptions.GET("/me", handleGetMySubscription)
			subscriptions.DELETE("/me", handleDeleteMySubscription)
		}
		v1.POST("/payments/webhook", handlePaymentsWebhook)

		// 16) Admin & Moderation
		admin := v1.Group("/admin")
		{
			admin.GET("/overview", handleAdminOverview)
			admin.GET("/flags", handleAdminGetFlags)
			admin.PATCH("/flags/:id/resolve", handleAdminResolveFlag)
			admin.GET("/auditlogs", handleAdminGetAuditLogs)
			admin.GET("/health", handleAdminHealth)
		}

		// 17) Notifications (server-side events)
		realtime := v1.Group("/realtime")
		{
			realtime.GET("/stream", handleRealtimeStream)
		}

		// 18) Localization
		i18n := v1.Group("/i18n")
		{
			i18n.GET("/locales", handleI18nGetLocales)
			i18n.GET("/strings", handleI18nGetStrings)
		}

		// 19) Analytics
		analytics := v1.Group("/analytics")
		{
			analytics.POST("/events", handleAnalyticsEvents)
		}
	}

	return router
}

func handleVerifyEmail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"verified":     true,
		"token":        "jwt.access.token.string",
		"refreshToken": "jwt.refresh.token.string",
		"user": gin.H{
			"id":   "user_123",
			"name": "Test User",
		},
	})
}

func handleForgot(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func handleReset(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 2) Users & Profiles
func handleGetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":     c.Param("id"),
		"name":   "Public User Name",
		"orgId":  "org_789",
		"badges": []string{"Top Contributor", "Verified"},
	})
}

func handleGetUserPreferences(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"notifications": gin.H{
			"email": true,
			"push":  false,
		},
		"preferredLang": "en",
	})
}

func handleUpdateUserPreferences(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"notifications": gin.H{
			"email": false,
			"push":  true,
		},
		"preferredLang": "am",
	})
}

func handleGetUserSummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"proceduresTracked": 5,
		"documents":         12,
		"remindersActive":   3,
	})
}

// 3) Organizations
func handleCreateOrg(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"organization": gin.H{
			"id":   "org_new_123",
			"name": "New Government Office",
			"type": "government",
		},
	})
}

func handleGetPendingOrgs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":   "org_pending_1",
				"name": "Pending Org A",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleApproveOrg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":       c.Param("id"),
		"approved": true,
	})
}

func handleGetOrgs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":   "org_456",
				"name": "Ministry of Innovation",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleGetOrg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":          c.Param("id"),
		"description": "Official government body for innovation.",
		"location":    "Addis Ababa",
		"type":        "government",
		"contact_info": gin.H{
			"email": "contact@moi.gov.et",
		},
		"phone_numbers": []string{"+25111223344"},
	})
}

func handleUpdateOrg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":          c.Param("id"),
		"description": "Updated description.",
	})
}

func handleGetOrgFeedback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":   "fb_1",
				"body": "This was very helpful!",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

// 4) Categories & Taxonomy

func handleUpdateCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "Updated Category Name"})
}

// 5) Procedures (core)

func handleGetProcedures(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":    "prc_123",
				"title": "Passport Renewal",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleGetProcedure(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":      "prc_123",
		"orgId":   "org_456",
		"title":   "Passport Renewal",
		"slug":    "passport-renewal",
		"summary": "Renew your Ethiopian passport in 5 steps.",
		"requirements": []gin.H{
			{"text": "2 passport photos"},
			{"text": "Old passport"},
		},
		"steps": []gin.H{
			{"order": 1, "text": "Book appointment"},
			{"order": 2, "text": "Submit documents"},
		},
		"fees": []gin.H{
			{"label": "Processing", "amount": 500, "currency": "ETB"},
		},
		"processingTime": gin.H{"minDays": 7, "maxDays": 14},
		"offices": []gin.H{
			{"city": "Addis Ababa", "address": "...", "hours": "Mon–Fri"},
		},
		"documentsRequired": []gin.H{
			{"name": "Application Form", "templateUrl": nil},
		},
		"tags":             []string{"passport", "id"},
		"languageVersions": gin.H{"enId": "prc_123", "amId": "prc_789"},
		"verified":         true,
		"updatedAt":        "2025-08-20T12:00:00Z",
	})
}

func handleUpdateProcedure(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "title": "Updated Procedure Title"})
}

func handleDeleteProcedure(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func handleVerifyProcedure(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "verified": true})
}

func handleGetProcedureAudit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"user":      "admin_user",
				"change":    "Set verified to true",
			},
		},
	})
}

func handleGetPopularProcedures(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{{"id": "prc_123", "title": "Passport Renewal", "views": 1050}},
	})
}

func handleGetRecentProcedures(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{{"id": "prc_789", "title": "New Business License", "updatedAt": time.Now().UTC().Format(time.RFC3339)}},
	})
}

// 6) Search & Discovery
func handleSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"procedures": []gin.H{{"id": "prc_123", "title": "Passport Renewal"}},
		"orgs":       []gin.H{{"id": "org_456", "name": "Ministry of Passport Services"}},
		"posts":      []gin.H{},
		"notices":    []gin.H{},
	})
}

// 7) Checklists & Progress
func handleCreateChecklist(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"id":          "chk_123",
		"userId":      "user_abc",
		"procedureId": "prc_123",
		"checklists": []gin.H{
			{"checklistID": "item_1", "text": "2 passport photos", "done": false},
			{"checklistID": "item_2", "text": "Old passport", "done": false},
		},
		"status":    "NOT_STARTED",
		"percent":   0,
		"updatedAt": time.Now().UTC().Format(time.RFC3339),
	})
}

func handleGetMyProcedures(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"userProcedureId": "up_123",
				"procedureId":     "prc_123",
				"procedureTitle":  "Passport Renewal",
				"status":          "IN_PROGRESS",
				"percent":         50,
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleGetChecklistByUserProcedureId(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":          c.Param("userProcedureId"),
		"userId":      "user_abc",
		"procedureId": "prc_123",
		"checklists": []gin.H{
			{"checklistID": "item_1", "text": "2 passport photos", "done": true},
			{"checklistID": "item_2", "text": "Old passport", "done": false},
		},
		"status":    "IN_PROGRESS",
		"percent":   50,
		"updatedAt": time.Now().UTC().Format(time.RFC3339),
	})
}

func handleUpdateChecklistItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"checklistID": c.Param("checklistID"), "done": true})
}

// 8) Documents & File Vault
func handleUploadSignature(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"signature": "cloudinary_signature_string",
		"timestamp": time.Now().Unix(),
		"apiKey":    "cloudinary_api_key",
		"cloudName": "your_cloud_name",
	})
}

func handleCreateDocument(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"id":      "doc_123",
		"userId":  "user_abc",
		"name":    "My Passport Scan",
		"fileUrl": "https://res.cloudinary.com/...",
		"type":    "passport",
	})
}

func handleGetDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":        "doc_123",
				"name":      "My Passport Scan",
				"fileUrl":   "https://res.cloudinary.com/...",
				"type":      "passport",
				"expiresOn": time.Now().AddDate(5, 0, 0).UTC().Format(time.RFC3339),
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleGetDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":        "doc_123",
		"userId":    "user_abc",
		"name":      "My Passport Scan",
		"fileUrl":   "https://res.cloudinary.com/...",
		"type":      "passport",
		"tags":      []string{"prc_123"},
		"issuedOn":  time.Now().AddDate(-5, 0, 0).UTC().Format(time.RFC3339),
		"expiresOn": time.Now().AddDate(5, 0, 0).UTC().Format(time.RFC3339),
		"ocrData":   gin.H{"name": "Test User"},
		"size":      1024,
	})
}

func handleUpdateDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "Updated Document Name"})
}

func handleDeleteDocument(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 9) Reminders & Notifications
func handleCreateReminder(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "rem_123", "title": "Renew Passport", "status": "ACTIVE"})
}

func handleGetReminders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":      "rem_123",
				"userId":  "user_abc",
				"title":   "Renew Passport",
				"dueAt":   time.Now().AddDate(0, 6, 0).UTC().Format(time.RFC3339),
				"channel": "email",
				"status":  "ACTIVE",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleUpdateReminder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "CANCELLED"})
}

func handleDeleteReminder(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func handleGetNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":   "notif_1",
				"body": "Your procedure 'Passport Renewal' has been updated.",
				"read": false,
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleMarkNotificationRead(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 10) Discussions (Community)
func handleUpvoteDiscussion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "votes": 13})
}

func handleDownvoteDiscussion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "votes": 11})
}

func handleReportDiscussion(c *gin.Context) {
	c.Status(http.StatusAccepted)
}

// Feedback
func handleUpvoteFeedback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "votes": 6})
}

// 12) Official Notices
func handleCreateNotice(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "ntc_new", "title": "New Official Notice"})
}

func handleGetNotices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":    "ntc_1",
				"title": "Holiday Office Closure",
			},
		},
		"page":    1,
		"limit":   20,
		"total":   1,
		"hasNext": false,
	})
}

func handleGetNotice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":            "ntc_1",
		"orgId":         "org_456",
		"title":         "Holiday Office Closure",
		"body":          "All offices will be closed on...",
		"pinned":        true,
		"effectiveFrom": "2025-09-10T00:00:00Z",
		"createdAt":     "2025-08-15T10:00:00Z",
	})
}

func handleUpdateNotice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "pinned": false})
}

func handleDeleteNotice(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// 13) AI Guidance
func handleAIGuide(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"response":  "Here are the verified steps...",
		"source":    "official",
		"citations": []gin.H{{"type": "procedure", "id": "prc_123"}},
		"verified":  true,
	})
}

func handleAIGetHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"history": []gin.H{
			{
				"id":       "id",
				"request":  "How to renew passport?",
				"response": "Here are the verified steps...",
				"source":   "official",
				"procedures": []gin.H{{
					"id":   "id",
					"name": "name",
				},
					{
						"id":   "id2",
						"name": "name2",
					},
				},
			},
		},
	})
}

func handleAIMarkNotVerified(c *gin.Context) {
	c.Status(http.StatusAccepted)
}

func handleAISpeechToText(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"text": "This is the recognized text from speech."})
}

// 14) Direct Messages
func handleCreateDMThread(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "dm_thread_1", "status": "open"})
}

func handleGetDMThreads(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":            "dm_thread_1",
				"orgId":         "org_456",
				"status":        "open",
				"lastMessageAt": time.Now().UTC().Format(time.RFC3339),
			},
		},
	})
}

func handleGetDMThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": c.Param("id"),
		"messages": []gin.H{
			{
				"id":   "msg_1",
				"from": "user",
				"body": "Hello, I have a question.",
			},
		},
	})
}

func handleCreateDMMessage(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": "msg_new", "threadId": c.Param("id"), "body": "This is a new message."})
}

func handleCloseDMThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "closed"})
}

// 15) Subscriptions & Payments
func handleGetPlans(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":           "plan_pro",
				"name":         "Pro",
				"priceMonthly": 100,
				"currency":     "ETB",
				"features":     []string{"Direct Messages", "Auto-tick"},
			},
		},
	})
}

func handleCreateSubscription(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pending", "clientSecret": "stripe_client_secret_string"})
}

func handleGetMySubscription(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"planId":     "plan_pro",
		"status":     "active",
		"currentEnd": time.Now().AddDate(0, 1, 0).UTC().Format(time.RFC3339),
	})
}

func handleDeleteMySubscription(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "cancelled", "message": "Subscription will be cancelled at the end of the current period."})
}

func handlePaymentsWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"received": true})
}

// 16) Admin & Moderation
func handleAdminOverview(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"uptime":        "99.99%",
		"dailyActives":  1500,
		"contentCounts": gin.H{"procedures": 120, "orgs": 30},
	})
}

func handleAdminGetFlags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"id":      "flag_1",
				"content": "Inappropriate discussion post.",
				"type":    "discussion",
				"status":  "pending",
			},
		},
	})
}

func handleAdminResolveFlag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "resolved"})
}

func handleAdminGetAuditLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": []gin.H{{"timestamp": time.Now(), "actor": "admin", "action": "deleted procedure prc_xyz"}}})
}

func handleAdminHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"services": gin.H{
			"database":   "connected",
			"cache":      "connected",
			"gemini_api": "ok",
		},
	})
}

// 17) Notifications (server-side events)
func handleRealtimeStream(c *gin.Context) {
	c.String(http.StatusOK, "This would be an SSE stream.")
}

// 18) Localization
func handleI18nGetLocales(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"supported": []string{"en", "am"}})
}

func handleI18nGetStrings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"welcomeMessage": "Welcome to EthioGuide",
		"search":         "Search",
	})
}

// 19) Analytics
func handleAnalyticsEvents(c *gin.Context) {
	c.Status(http.StatusAccepted)
}
