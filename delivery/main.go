package main

import (
	"EthioGuide/config"
	"EthioGuide/delivery/controller"
	"EthioGuide/delivery/router"
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
	// FIX 1: Initialize the TokenRepository, as it's a required dependency for UserUsecase.
	tokenRepo := repository.NewTokenRepository(db)
	checklistRepo := repository.NewChecklistRepository(db)

	// --- Infrastructure Services ---
	// These are concrete implementations of external services.
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	aiService, err := infrastructure.NewGeminiAIService(cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		log.Printf("WARN: Failed to initialize AI service: %v. AI features will be unavailable.", err)
	}

	// --- Use Cases ---
	// Use cases orchestrate the business logic, using repositories and services.
	// FIX 2: Pass all the required dependencies to the NewUserUsecase constructor.
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		tokenRepo, // Added the missing token repository
		passwordService,
		jwtService,
		cfg.UsecaseTimeout,
	)
	geminiUsecase := usecase.NewGeminiUsecase(aiService, cfg.UsecaseTimeout) // Reduced timeout for consistency
	checklistUsecase := usecase.NewChecklistUsecase(checklistRepo)
	// --- Controllers ---
	// Controllers handle the HTTP layer, delegating logic to use cases.
	userController := controller.NewUserController(userUsecase, checklistUsecase, cfg.JWTRefreshTTL)
	geminiController := controller.NewGeminiController(geminiUsecase)

	// --- Middleware ---
	// Middleware is created to be injected into the router.
	authMiddleware := infrastructure.AuthMiddleware(jwtService)
	proOnlyMiddleware := infrastructure.ProOnlyMiddleware()
	requireAdminRole := infrastructure.RequireRole(domain.RoleAdmin)

	// --- Router Setup ---
	// The router is configured with all the controllers and middleware.
	appRouter := router.SetupRouter(
		userController,
		geminiController,
		authMiddleware,
		proOnlyMiddleware,
		requireAdminRole,
	)

	// --- Start Server ---
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := appRouter.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
