package main

import (
	"EthioGuide/config"
	"EthioGuide/infrastructure"
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
	defer client.Disconnect(context.Background())
	db := client.Database(cfg.DBName)
	log.Println("MongoDB connected.")
	println("Ethio Guide")

	// --- Infrastructure Services ---
	// Pass values from the cfg struct to the service constructors.
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	emailService := infrastructure.NewSMTPEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	aiService, err := infrastructure.NewGeminiAIService(cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		log.Printf("WARN: Failed to initialize AI service: %v. AI features will be unavailable.", err)
	}
	googleOAuth2Service, err := infrastructure.NewGoogleOAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)
	if err != nil {
		log.Println("WARN: Google OAuth credentials are not set. Sign in with Google will fail.", err)
	}
	redisService, err := infrastructure.NewRedisService(context.Background(), cfg.RedisUrl, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("FATAL: Redis connection failed.Error: %v", err)
	}
	defer redisService.Close()
	rateLimiter := infrastructure.NewRateLimiter(redisService)
	cacheService := infrastructure.NewRedisCacheService(redisService)

	println(db, passwordService, jwtService, emailService, aiService, googleOAuth2Service, rateLimiter, cacheService)
}
