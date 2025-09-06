package main

import (
	"EthioGuide/config"
	"EthioGuide/delivery/controller"
	"EthioGuide/delivery/router"
	_ "EthioGuide/docs"
	"EthioGuide/domain"
	"EthioGuide/infrastructure"
	"EthioGuide/repository"
	"EthioGuide/usecase"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           EthioGuide API
// @version         1.0
// @description     This is the API server for the EthioGuide application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      ethio-guide-backend.onrender.com
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and a JWT token."

func main() {
	// --- Load Configuration ---
	cfg := config.Load()

	// --- MongoDB Setup ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	db := client.Database(cfg.DBName)
	log.Println("Successfully connected to MongoDB.")

	// --- Repositories ---
	// Repositories are the first layer to be initialized as they only depend on the database.
	userRepo := repository.NewAccountRepository(db)
	procedureRepo := repository.NewProcedureRepository(db)
	preferencesRepo := repository.NewPreferencesRepository(db)
	catagoryRepo := repository.NewCategoryRepository(db, "catagories")
	tokenRepo := repository.NewTokenRepository(db)
	feedbackRepo := repository.NewFeedbackRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	postRepo := repository.NewPostRepository(db)
	searchRepo := repository.NewSearchRepository(db)
	checklistRepo := repository.NewChecklistRepository(db)
	aiChatRepo := repository.NewAIChatRepository(db)

	// --- Infrastructure Services ---
	// These are concrete implementations of external services.
	emailservice := infrastructure.NewSMTPEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, cfg.VerificationFrontendUrl, cfg.ResetPasswordFrontendUrl)
	passwordService := infrastructure.NewPasswordService()
	googleService, err := infrastructure.NewGoogleOAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)
	if err != nil {
		log.Printf("WARN: Failed to initialize Google Oaut service: %v. Google Sign in will be unavailable.", err)
	}
	jwtService := infrastructure.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTL, cfg.JWTRefreshTTL, cfg.JWTUtilityTTL)
	aiService, err := infrastructure.NewGeminiAIService(cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		log.Printf("WARN: Failed to initialize AI service: %v. AI features will be unavailable.", err)
	}
	var apiKeys []string
	apiKeys = append(apiKeys, cfg.EmbeddingApiKey)
	embeddingService, err := infrastructure.NewEmbeddingService(apiKeys, cfg.EmbeddingUrl)
	if err != nil {
		log.Printf("WARN: Failed to initialize the Embedding service: %v. Embedding service will be unavailable.", err)
	}

	// --- Use Cases ---
	// Use cases orchestrate the business logic, using repositories and services.
	// FIX 2: Pass all the required dependencies to the NewUserUsecase constructor.
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		tokenRepo, // Added the missing token repository
		passwordService,
		jwtService,
		googleService,
		emailservice,
		cfg.UsecaseTimeout,
	)
	procedureUsecase := usecase.NewProcedureUsecase(procedureRepo, embeddingService, cfg.UsecaseTimeout)
	catagoryUsecase := usecase.NewCategoryUsecase(catagoryRepo, cfg.UsecaseTimeout)
	geminiUsecase := usecase.NewGeminiUsecase(aiService, cfg.UsecaseTimeout) // Reduced timeout for consistency
	feedbackUsecase := usecase.NewFeedbackUsecase(feedbackRepo, procedureRepo, cfg.UsecaseTimeout)
	noticeUsecase := usecase.NewNoticeUsecase(noticeRepo)
	preferencesUsecase := usecase.NewPreferencesUsecase(preferencesRepo)
	aiChatUsecase := usecase.NewChatUsecase(embeddingService, procedureRepo, aiChatRepo, aiService)

	postUsecase := usecase.NewPostUseCase(postRepo, cfg.UsecaseTimeout)
	searchUsecase := usecase.NewSearchUsecase(searchRepo, cfg.UsecaseTimeout)
	checklistUsecase := usecase.NewChecklistUsecase(checklistRepo)
	// --- Controllers ---
	// Controllers handle the HTTP layer, delegating logic to use cases.
	userController := controller.NewUserController(userUsecase, searchUsecase, checklistUsecase, cfg.JWTRefreshTTL)
	procedureController := controller.NewProcedureController(procedureUsecase)
	catagoryController := controller.NewCategoryController(catagoryUsecase)
	geminiController := controller.NewGeminiController(geminiUsecase)
	feedbackController := controller.NewFeedbackController(feedbackUsecase)
	noticeController := controller.NewNoticeController(noticeUsecase)
	preferencesController := controller.NewPreferencesController(preferencesUsecase)
	aiChatController := controller.NewAIChatController(aiChatUsecase)

	postController := controller.NewPostController(postUsecase)
	// --- Middleware ---
	// Middleware is created to be injected into the router.
	authMiddleware := infrastructure.AuthMiddleware(jwtService)
	proOnlyMiddleware := infrastructure.ProOnlyMiddleware()
	requireAdminRole := infrastructure.RequireRole(domain.RoleAdmin)
	requireAdminOrOrgRole := infrastructure.RequireRole(domain.RoleAdmin, domain.RoleOrg)

	// --- Router Setup ---
	// The router is configured with all the controllers and middleware.
	appRouter := router.SetupRouter(
		userController,
		procedureController,
		catagoryController,
		geminiController,
		feedbackController,
		postController,
		noticeController,
		preferencesController,
		aiChatController,
		authMiddleware,
		proOnlyMiddleware,
		requireAdminRole,
		requireAdminOrOrgRole,
	)

	// --- Start Server ---
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := appRouter.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
