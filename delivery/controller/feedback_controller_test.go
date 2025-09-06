package controller_test

// import (
// 	. "EthioGuide/delivery/controller"
// 	"EthioGuide/domain"
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// // --- Mock for IFeedbackUsecase ---
// type MockFeedbackUsecase struct {
//     mock.Mock
// }

// func (m *MockFeedbackUsecase) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
//     args := m.Called(ctx, feedback)
//     return args.Error(0)
// }

// func (m *MockFeedbackUsecase) GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
//     args := m.Called(ctx, procedureID, filter)
//     var feedbacks []*domain.Feedback
//     if args.Get(0) != nil {
//         feedbacks = args.Get(0).([]*domain.Feedback)
//     }
//     return feedbacks, args.Get(1).(int64), args.Error(2)
// }

// func (m *MockFeedbackUsecase) UpdateFeedbackStatus(ctx context.Context, feedbackID, userID string, status domain.FeedbackStatus, adminResponse *string) error {
//     args := m.Called(ctx, feedbackID, userID, status, adminResponse)
//     return args.Error(0)
// }
// func (m *MockFeedbackUsecase) GetAllFeedbacks(ctx context.Context, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
//     args := m.Called(ctx, filter)
//     return args.Get(0).([]*domain.Feedback), args.Get(1).(int64), args.Error(2)
// }

// // --- Test Suite ---
// type FeedbackControllerTestSuite struct {
//     suite.Suite
//     router     *gin.Engine
//     mockUsecase *MockFeedbackUsecase
//     controller *FeedbackController
//     recorder   *httptest.ResponseRecorder
// }

// func (s *FeedbackControllerTestSuite) SetupSuite() {
//     gin.SetMode(gin.TestMode)
// }

// func (s *FeedbackControllerTestSuite) SetupTest() {
//     s.recorder = httptest.NewRecorder()
//     s.router = gin.Default()
//     s.mockUsecase = new(MockFeedbackUsecase)
//     s.controller = NewFeedbackController(s.mockUsecase)

//     // Setup routes
//     s.router.POST("/procedures/:id/feedback", func(c *gin.Context) {
//         // Simulate authenticated user
//         c.Set("userID", "user123")
//         s.controller.SubmitFeedback(c)
//     })
//     s.router.GET("/procedures/:id/feedback", s.controller.GetAllFeedbacksForProcedure)
//     s.router.PATCH("/feedback/:id", func(c *gin.Context) {
//         c.Set("userID", "org456")
//         s.controller.UpdateFeedbackStatus(c)
//     })
// }

// func TestFeedbackControllerTestSuite(t *testing.T) {
//     suite.Run(t, new(FeedbackControllerTestSuite))
// }

// // --- Tests ---

// func (s *FeedbackControllerTestSuite) TestSubmitFeedback() {
//     reqBody := FeedbackCreateRequest{
//         Content: "Great info!",
//         Type:    "thanks",
//         Tags:    []string{"tag1"},
//     }
//     jsonBody, _ := json.Marshal(reqBody)
//     // feedback := &domain.Feedback{
//     //     UserID:      "user123",
//     //     ProcedureID: "proc789",
//     //     Content:     reqBody.Content,
//     //     Type:        domain.FeedbackType(reqBody.Type),
//     //     Tags:        reqBody.Tags,
//     // }

//     s.Run("Success", func() {
//         s.SetupTest()
//         s.mockUsecase.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*domain.Feedback")).Return(nil).Once()

//         req, _ := http.NewRequest(http.MethodPost, "/procedures/proc789/feedback", bytes.NewBuffer(jsonBody))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusCreated, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), "Feedback submitted successfully")
//         s.mockUsecase.AssertExpectations(s.T())
//     })

//     s.Run("Failure - Invalid Body", func() {
//         s.SetupTest()
//         req, _ := http.NewRequest(http.MethodPost, "/procedures/proc789/feedback", bytes.NewBuffer([]byte(`{invalid json`)))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusBadRequest, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), "Invalid request body")
//     })

//     s.Run("Failure - Usecase Error", func() {
//         s.SetupTest()
//         s.mockUsecase.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*domain.Feedback")).Return(domain.ErrConflict).Once()

//         req, _ := http.NewRequest(http.MethodPost, "/procedures/proc789/feedback", bytes.NewBuffer(jsonBody))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusConflict, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), domain.ErrConflict.Error())
//         s.mockUsecase.AssertExpectations(s.T())
//     })
// }

// func (s *FeedbackControllerTestSuite) TestGetAllFeedbacksForProcedure() {
//     // filter := &domain.FeedbackFilter{Page: 1, Limit: 10}
//     feedbacks := []*domain.Feedback{
//         {
//             ID:          "fb1",
//             UserID:      "user123",
//             ProcedureID: "proc789",
//             Content:     "Great info!",
//             Type:        domain.ThanksFeedback,
//             Status:      domain.NewFeedback,
//             CreatedAT:   time.Now(),
//             UpdatedAT:   time.Now(),
//         },
//     }
//     s.Run("Success", func() {
//         s.SetupTest()
//         s.mockUsecase.On("GetAllFeedbacksForProcedure", mock.Anything, "proc789", mock.AnythingOfType("*domain.FeedbackFilter")).
//             Return(feedbacks, int64(1), nil).Once()

//         req, _ := http.NewRequest(http.MethodGet, "/procedures/proc789/feedback?page=1&limit=10", nil)
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusOK, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), "feedbacks")
//         s.mockUsecase.AssertExpectations(s.T())
//     })

//     s.Run("Failure - Usecase Error", func() {
//         s.SetupTest()
//         s.mockUsecase.On("GetAllFeedbacksForProcedure", mock.Anything, "proc789", mock.AnythingOfType("*domain.FeedbackFilter")).
//             Return(nil, int64(0), domain.ErrNotFound).Once()

//         req, _ := http.NewRequest(http.MethodGet, "/procedures/proc789/feedback?page=1&limit=10", nil)
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusNotFound, s.recorder.Code)
//         s.mockUsecase.AssertExpectations(s.T())
//     })
// }

// func (s *FeedbackControllerTestSuite) TestUpdateFeedbackStatus() {
//     reqBody := FeedbackStatePatchRequest{
//         Status:        "resolved",
//         AdminResponse: "Fixed!",
//     }
//     jsonBody, _ := json.Marshal(reqBody)

//     s.Run("Success", func() {
//         s.SetupTest()
//         s.mockUsecase.On("UpdateFeedbackStatus", mock.Anything, "fb1", "org456", domain.ResolvedFeedback, mock.AnythingOfType("*string")).Return(nil).Once()

//         req, _ := http.NewRequest(http.MethodPatch, "/feedback/fb1", bytes.NewBuffer(jsonBody))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusOK, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), "Feedback status updated successfully")
//         s.mockUsecase.AssertExpectations(s.T())
//     })

//     s.Run("Failure - Invalid Body", func() {
//         s.SetupTest()
//         req, _ := http.NewRequest(http.MethodPatch, "/feedback/fb1", bytes.NewBuffer([]byte(`{invalid json`)))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusBadRequest, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), "Invalid request body")
//     })

//     s.Run("Failure - Usecase Error", func() {
//         s.SetupTest()
//         s.mockUsecase.On("UpdateFeedbackStatus", mock.Anything, "fb1", "org456", domain.ResolvedFeedback, mock.AnythingOfType("*string")).Return(domain.ErrPermissionDenied).Once()

//         req, _ := http.NewRequest(http.MethodPatch, "/feedback/fb1", bytes.NewBuffer(jsonBody))
//         req.Header.Set("Content-Type", "application/json")
//         s.router.ServeHTTP(s.recorder, req)

//         s.Equal(http.StatusForbidden, s.recorder.Code)
//         s.Contains(s.recorder.Body.String(), domain.ErrPermissionDenied.Error())
//         s.mockUsecase.AssertExpectations(s.T())
//     })
// }
