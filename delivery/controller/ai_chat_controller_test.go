package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"EthioGuide/delivery/controller"
	// "EthioGuide/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Usecase ---
type MockAIChatUsecase struct {
    mock.Mock
}

func (m *MockAIChatUsecase) AIchat(ctx context.Context, query string) (string, error) {
    args := m.Called(ctx, query)
    return args.String(0), args.Error(1)
}

// --- Test Suite ---
type AIChatControllerSuite struct {
    suite.Suite
    router   *gin.Engine
    mockUC   *MockAIChatUsecase
    recorder *httptest.ResponseRecorder
}

func (s *AIChatControllerSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    s.recorder = httptest.NewRecorder()
    s.mockUC = new(MockAIChatUsecase)
    ctl := controller.NewAIChatController(s.mockUC)

    s.router = gin.Default()
    s.router.POST("/ai-chat", ctl.AIChatController)
}

func TestAIChatControllerSuite(t *testing.T) {
    suite.Run(t, new(AIChatControllerSuite))
}

func (s *AIChatControllerSuite) TestAIChat_Success() {
    reqBody := map[string]interface{}{"query": "Hello AI"}
    body, _ := json.Marshal(reqBody)

    // Set up mock expectation
    s.mockUC.On("AIchat", mock.Anything, "Hello AI").Return("Hi there!", nil).Once()

    req, _ := http.NewRequest(http.MethodPost, "/ai-chat", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    s.router.ServeHTTP(s.recorder, req)

    s.Equal(http.StatusOK, s.recorder.Code)
    var resp map[string]interface{}
    json.Unmarshal(s.recorder.Body.Bytes(), &resp)
    s.Equal("Hi there!", resp["answer"])
    s.mockUC.AssertExpectations(s.T())
}

func (s *AIChatControllerSuite) TestAIChat_BadRequest() {
    reqBody := map[string]interface{}{"wrongField": "Hello AI"}
    body, _ := json.Marshal(reqBody)

    req, _ := http.NewRequest(http.MethodPost, "/ai-chat", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    s.router.ServeHTTP(s.recorder, req)

    s.Equal(http.StatusBadRequest, s.recorder.Code)
    var resp map[string]interface{}
    json.Unmarshal(s.recorder.Body.Bytes(), &resp)
    s.Equal("Invalid request", resp["error"])
}

func (s *AIChatControllerSuite) TestAIChat_ErrorFromUsecase() {
    reqBody := map[string]interface{}{"query": "fail"}
    body, _ := json.Marshal(reqBody)

    s.mockUC.On("AIchat", mock.Anything, "fail").Return("", assert.AnError).Once()

    req, _ := http.NewRequest(http.MethodPost, "/ai-chat", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    s.router.ServeHTTP(s.recorder, req)

    s.Equal(http.StatusInternalServerError, s.recorder.Code)
    var resp map[string]interface{}
    json.Unmarshal(s.recorder.Body.Bytes(), &resp)
    s.Contains(resp["error"], assert.AnError.Error())
}