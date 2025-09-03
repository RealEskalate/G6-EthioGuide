package infrastructure_test

import (
	"EthioGuide/domain" // Your actual domain package
	. "EthioGuide/infrastructure"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Setup ---
// This section creates a mock JWTService and placeholder domain types
// so we can test the middleware in isolation.

// MockJWTService is a mock implementation of IJWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID string, role domain.Role) (string, *domain.JWTClaims, error) {
	// m.Called() tells the mock object that this method has been called.
	// It records the call and its arguments.
	args := m.Called(userID, role)

	// We retrieve the return values that we configured in our test.
	var claims *domain.JWTClaims
	if args.Get(1) != nil {
		claims = args.Get(1).(*domain.JWTClaims)
	}

	return args.String(0), claims, args.Error(2)
}

func (m *MockJWTService) GenerateRefreshToken(userID string) (string, *domain.JWTClaims, error) {
	args := m.Called(userID)

	var claims *domain.JWTClaims
	if args.Get(1) != nil {
		claims = args.Get(1).(*domain.JWTClaims)
	}

	return args.String(0), claims, args.Error(2)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)

	// Safely handle the possibility of a nil claims return
	var claims *domain.JWTClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*domain.JWTClaims)
	}

	return claims, args.Error(1)
}

func (m *MockJWTService) ParseExpiredToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)

	var claims *domain.JWTClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*domain.JWTClaims)
	}

	return claims, args.Error(1)
}

func (m *MockJWTService) GetRefreshTokenExpiry() time.Duration {
	args := m.Called()

	// The return value is retrieved by its index.
	// We need to perform a type assertion to the expected type.
	return args.Get(0).(time.Duration)
}

// --- Test Suite Definition ---

type MiddlewareTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mockJWTService *MockJWTService
	recorder       *httptest.ResponseRecorder
}

// SetupTest runs before each test in the suite
func (s *MiddlewareTestSuite) SetupTest() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockJWTService = new(MockJWTService)
}

// TestMiddlewareTestSuite is the entry point for the test suite
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

// --- AuthMiddleware Tests ---

func (s *MiddlewareTestSuite) TestAuthMiddleware_NoAuthHeader() {
	// Setup the route with the middleware
	s.router.GET("/test", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Create a request without the Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	s.router.ServeHTTP(s.recorder, req)

	// Assertions
	s.Equal(http.StatusUnauthorized, s.recorder.Code)

	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("Authorization header is required", response["error"])

	// Ensure the JWT service was never called
	s.mockJWTService.AssertNotCalled(s.T(), "ValidateToken")
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_MalformedHeader() {
	s.router.GET("/test", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusUnauthorized, s.recorder.Code)
	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("Authorization header format must be Bearer {token}", response["error"])
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	invalidToken := "this.is.an.invalid.token"
	// Expect ValidateToken to be called with the token and return an error
	s.mockJWTService.On("ValidateToken", invalidToken).Return(nil, errors.New("token is invalid"))

	s.router.GET("/test", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusUnauthorized, s.recorder.Code)
	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("invalid or expired token", response["error"])

	// Verify that the mock expectation was met
	s.mockJWTService.AssertExpectations(s.T())
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_Success() {
	validToken := "a.valid.jwt.token"
	expectedClaims := &domain.JWTClaims{
		UserID:       "user-123",
		Role:         domain.RoleUser,
		Subscription: "pro",
	}
	s.mockJWTService.On("ValidateToken", validToken).Return(expectedClaims, nil)

	var capturedUserID, capturedUserRole, capturedUserSubscription interface{}

	s.router.GET("/test", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		// This handler runs only if middleware succeeds.
		// We capture the values set in the context.
		capturedUserID, _ = c.Get("userID")
		capturedUserRole, _ = c.Get("userRole")
		capturedUserSubscription, _ = c.Get("userSubscription")
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusOK, s.recorder.Code)
	s.mockJWTService.AssertExpectations(s.T())

	// Assert that the context values were set correctly
	s.Equal(expectedClaims.UserID, capturedUserID)
	s.Equal(expectedClaims.Role, capturedUserRole)
	s.Equal(expectedClaims.Subscription, capturedUserSubscription)
}

// --- ProOnlyMiddleware Tests ---

// Helper function to create a context with predefined values
func setContext(c *gin.Context, subscription string, role domain.Role) {
	c.Set("userSubscription", subscription)
	c.Set("userRole", role)
}

func (s *MiddlewareTestSuite) TestProOnlyMiddleware_Success() {
	s.router.GET("/pro", func(c *gin.Context) {
		setContext(c, "pro", domain.RoleUser)
	}, ProOnlyMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/pro", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusOK, s.recorder.Code)
}

func (s *MiddlewareTestSuite) TestProOnlyMiddleware_InsufficientSubscription() {
	s.router.GET("/pro", func(c *gin.Context) {
		setContext(c, "free", domain.RoleUser)
	}, ProOnlyMiddleware(), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/pro", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusForbidden, s.recorder.Code)
	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("Insufficient subscription permissions", response["error"])
}

func (s *MiddlewareTestSuite) TestProOnlyMiddleware_MissingSubscription() {
	// No context-setting middleware
	s.router.GET("/pro", ProOnlyMiddleware(), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/pro", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusForbidden, s.recorder.Code)
	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("Subscription information is missing", response["error"])
}

// --- RequireRole Tests ---

func (s *MiddlewareTestSuite) TestRequireRole_Success() {
	// User is an Admin, route requires Admin or Organization
	s.router.GET("/admin", func(c *gin.Context) {
		setContext(c, "pro", domain.RoleAdmin)
	}, RequireRole(domain.RoleAdmin, domain.RoleOrg), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusOK, s.recorder.Code)
}

func (s *MiddlewareTestSuite) TestRequireRole_InsufficientRole() {
	// User is a User, route requires Admin
	s.router.GET("/admin", func(c *gin.Context) {
		setContext(c, "pro", domain.RoleUser)
	}, RequireRole(domain.RoleAdmin), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusForbidden, s.recorder.Code)
	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal("Insufficient role permissions", response["error"])
}

func (s *MiddlewareTestSuite) TestRequireRole_MissingRole() {
	// No context-setting middleware
	s.router.GET("/admin", RequireRole(domain.RoleAdmin), func(c *gin.Context) {})

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusForbidden, s.recorder.Code)
	// Note: Your original code returns a string, not JSON, for this specific error.
	s.Contains(s.recorder.Body.String(), "Role information is missing")
}
