package infrastructure_test

import (
	"EthioGuide/domain"
	. "EthioGuide/infrastructure"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupService is a helper to create a new service for tests.
func setupService() (domain.IJWTService, string, string) {
	secret := "my-super-secret-key-for-testing"
	issuer := "test-issuer"
	accessTTL := 15 * time.Minute
	refreshTTL := 24 * time.Hour
	utilityTTL := 1 * time.Hour
	jwtService := NewJWTService(secret, issuer, accessTTL, refreshTTL, utilityTTL)
	return jwtService, secret, issuer
}

func TestJWTService_GenerateAccessToken(t *testing.T) {
	jwtService, _, issuer := setupService()
	userID := "user-123"
	userRole := domain.RoleAdmin

	// Act
	tokenString, claims, err := jwtService.GenerateAccessToken(userID, userRole)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, userRole, claims.Role)
	assert.Equal(t, issuer, claims.Issuer)
	assert.NotEmpty(t, claims.ID, "JTI (ID) should not be empty")
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), claims.ExpiresAt.Time, 1*time.Second)
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	jwtService, _, issuer := setupService()
	userID := "user-456"

	// Act
	tokenString, claims, err := jwtService.GenerateRefreshToken(userID)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Empty(t, claims.Role, "Refresh token should not contain a role claim")
	assert.NotEmpty(t, claims.ID, "JTI (ID) should not be empty")
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, 1*time.Second)
}

func TestJWTService_Validation(t *testing.T) {
	jwtService, secret, _ := setupService()
	otherService, _, _ := setupService() // Just to have another instance for testing
	wrongSecretService := NewJWTService("a-different-secret", "test-issuer", 15*time.Minute, 24*time.Hour, 1*time.Hour)

	userID := "user-789"
	userRole := domain.RoleUser

	t.Run("Success - Valid token", func(t *testing.T) {
		tokenString, _, _ := jwtService.GenerateAccessToken(userID, userRole)
		claims, err := otherService.ValidateToken(tokenString)

		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	t.Run("Failure - Invalid signature (wrong secret)", func(t *testing.T) {
		tokenString, _, _ := jwtService.GenerateAccessToken(userID, userRole)
		_, err := wrongSecretService.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Failure - Token is expired", func(t *testing.T) {
		// Create a service that generates instantly expired tokens
		expiredService := NewJWTService(secret, "test-issuer", -5*time.Minute, -1*time.Hour, -1*time.Hour)
		tokenString, _, _ := expiredService.GenerateAccessToken(userID, userRole)

		_, err := jwtService.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims: token is expired")
	})

	t.Run("Failure - Malformed token string", func(t *testing.T) {
		// Arrange
		malformedToken := "this-is-not-a-jwt"

		// Act
		_, err := jwtService.ValidateToken(malformedToken)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is malformed")
	})
}

func TestJWTService_ParseExpiredToken(t *testing.T) {
	jwtService, secret, _ := setupService()
	userID := "user-exp"
	userRole := domain.RoleUser

	// Create a service that generates instantly expired tokens
	expiredService := NewJWTService(secret, "test-issuer", -5*time.Minute, 24*time.Hour, 1*time.Hour)

	t.Run("Success - Can parse an expired token", func(t *testing.T) {
		// Arrange
		tokenString, originalClaims, _ := expiredService.GenerateAccessToken(userID, userRole)

		// Act: Use the special parsing method
		parsedClaims, err := jwtService.ParseExpiredToken(tokenString)

		// Assert
		require.NoError(t, err, "ParseExpiredToken should not return an error for an expired token")
		require.NotNil(t, parsedClaims)

		assert.Equal(t, originalClaims.ID, parsedClaims.ID)
		assert.Equal(t, originalClaims.UserID, parsedClaims.UserID)
	})

	t.Run("Failure - Fails on invalid signature", func(t *testing.T) {
		// Arrange
		wrongSecretService := NewJWTService("wrong-secret", "test-issuer", -5*time.Minute, 24*time.Hour, 1*time.Hour)
		tokenString, _, _ := wrongSecretService.GenerateAccessToken(userID, userRole)

		// Act
		_, err := jwtService.ParseExpiredToken(tokenString)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})
}

func TestJWTService_GenerateUtilityToken(t *testing.T) {
	jwtService, _, issuer := setupService()
	userID := "user-utility"

	// Act
	tokenString, claims, err := jwtService.GenerateUtilityToken(userID)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Empty(t, claims.Role, "Utility token should not contain a role claim")
	assert.NotEmpty(t, claims.ID, "JTI (ID) should not be empty")
	assert.WithinDuration(t, time.Now().Add(1*time.Hour), claims.ExpiresAt.Time, 1*time.Second)
}
