package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
// Values are read from environment variables.
type Config struct {
	AppEnv         string
	ServerPort     string
	UsecaseTimeout time.Duration

	MongoURI string
	DBName   string

	RedisUrl      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret     string
	JWTIssuer     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	JWTUtilityTTL time.Duration

	GeminiAPIKey string
	GeminiModel  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	VerificationFrontendUrl  string
	ResetPasswordFrontendUrl string
	EmbeddingUrl             string
	EmbeddingApiKey          string
}

// Load loads the configuration from .env files and environment variables.
func Load() *Config {
	// Load .env file from current or parent directory
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, proceeding with environment defaults...")
		}
	}

	// Read and parse all environment variables
	accessTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TTL_MIN", "15"))
	refreshTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TTL_HR", "72"))
	utilityTTL, _ := strconv.Atoi(getEnv("JWT_UTILITY_TTL_HR", "72"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "2525"))
	frontendBaseUrl := getEnv("FRONTEND_BASE_URL", "http://localhost:3000")

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		ServerPort:     getEnv("PORT", "8080"),
		UsecaseTimeout: 5 * time.Second,

		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "ethio-guide-db"),

		RedisUrl:      getEnv("REDIS_URI", ""),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTSecret:     getEnv("JWT_SECRET", "a-very-secret-key-that-should-be-long-and-random"),
		JWTIssuer:     "ethio-guide-api",
		JWTAccessTTL:  time.Duration(accessTTL) * time.Minute,
		JWTRefreshTTL: time.Duration(refreshTTL) * time.Hour,
		JWTUtilityTTL: time.Duration(utilityTTL) * time.Hour,

		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.5-pro"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:  frontendBaseUrl + getEnv("GOOGLE_REDIRECT_URI", ""),

		SMTPHost: getEnv("SMTP_HOST", "smtp.mailtrap.io"),
		SMTPPort: smtpPort,
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom: getEnv("SMTP_FROM_EMAIL", "no-reply@example.com"),

		VerificationFrontendUrl:  frontendBaseUrl + getEnv("VERIFICATION_FRONTEND_URL", "/api/auth/verify"),
		ResetPasswordFrontendUrl: frontendBaseUrl + getEnv("RESET_PASSWORD_FRONTEND_URL", "/api/auth/forgot-password"),
		EmbeddingUrl:             getEnv("EMBEDDING_URL", "https://router.huggingface.co/hf-inference/models/sentence-transformers/all-MiniLM-L6-v2/pipeline/feature-extraction"),
		EmbeddingApiKey:          getEnv("HF_EMBEDDING_API_KEY", ""),
	}
}

// getEnv is a helper to read an environment variable or return a fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadForTest() *Config {
	// Load .env.test first for test-specific configurations.
	// We search in the current directory and the parent directory.
	godotenv.Load(".env.test")
	godotenv.Load("../.env.test")
	// Fallback to regular .env for shared configs.
	godotenv.Load(".env")
	godotenv.Load("../.env")

	// Use our test-specific helper for reading variables
	mongoURI := getTestEnv("MONGO_URI_TEST", "MONGO_URI", "mongodb://localhost:27017")
	dbName := getTestEnv("DB_NAME_TEST", "DB_NAME", "ethio-guide-db-test")

	// Perform the safety check here, within the config loader.
	if dbName == "ethio-guide-db" {
		log.Fatalf("FATAL: Cannot run tests on the main database '%s'. Set DB_NAME_TEST in your .env.test file.", dbName)
	}

	// Create a base config by calling the standard Load() function.
	// This ensures we get all the defaults and standard variables.
	cfg := Load()

	cfg.MongoURI = mongoURI
	cfg.DBName = dbName

	return cfg
}

// getTestEnv implements the desired fallback logic for environment variables.
// It checks for a test-specific key, then a normal key, and finally returns a default value.
// We make it unexported as it's a helper for this package.
func getTestEnv(testKey, fallbackKey, defaultValue string) string {
	if value, exists := os.LookupEnv(testKey); exists {
		return value
	}
	if value, exists := os.LookupEnv(fallbackKey); exists {
		return value
	}
	return defaultValue
}
