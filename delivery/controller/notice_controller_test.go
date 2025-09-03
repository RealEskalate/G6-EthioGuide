package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockNoticeUsecase struct{ mock.Mock }

var _ domain.INoticeUseCase = (*MockNoticeUsecase)(nil)

func (m *MockNoticeUsecase) CreateNotice(ctx context.Context, notice *domain.Notice) error {
	return m.Called(ctx, notice).Error(0)
}
func (m *MockNoticeUsecase) GetNoticesByFilter(ctx context.Context, filter *domain.NoticeFilter) ([]*domain.Notice, int64, error) {
	args := m.Called(ctx, filter)
	var res []*domain.Notice
	if args.Get(0) != nil {
		res = args.Get(0).([]*domain.Notice)
	}
	total := int64(0)
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return res, total, args.Error(2)
}
func (m *MockNoticeUsecase) UpdateNotice(ctx context.Context, id string, notice *domain.Notice) error {
	return m.Called(ctx, id, notice).Error(0)
}
func (m *MockNoticeUsecase) DeleteNotice(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type NoticeControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	ctrl   *NoticeController
	mockUC *MockNoticeUsecase
}

func (s *NoticeControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *NoticeControllerTestSuite) SetupTest() {
	s.router = gin.New()
	s.mockUC = new(MockNoticeUsecase)
	s.ctrl = NewNoticeController(s.mockUC)

	// Routes
	s.router.POST("/notices", s.ctrl.CreateNotice)
	s.router.GET("/notices", s.ctrl.GetNoticesByFilter)
	// Wrap Update/Delete to supply :id to controller method signatures
	s.router.PUT("/notices/:id", func(c *gin.Context) { s.ctrl.UpdateNotice(c, c.Param("id")) })
	s.router.DELETE("/notices/:id", func(c *gin.Context) { s.ctrl.DeleteNotice(c, c.Param("id")) })
}

func TestNoticeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(NoticeControllerTestSuite))
}

func (s *NoticeControllerTestSuite) TestCreateNotice() {
	body := NoticeDTO{
		OrganizationID: "64dbf7c7c12e2c3a4b5a6c7d",
		Title:          "Test Notice",
		Content:        "Body",
		Tags:           []string{"go"},
	}
	jsonBody, _ := json.Marshal(body)

	s.mockUC.On("CreateNotice", mock.Anything, mock.AnythingOfType("*domain.Notice")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/notices", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
	s.mockUC.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestGetNoticesByFilter() {
	expected := []*domain.Notice{{Title: "N1"}}
	total := int64(1)

	s.mockUC.
		On("GetNoticesByFilter", mock.Anything, mock.MatchedBy(func(f *domain.NoticeFilter) bool {
			// Expect controller to parse these correctly
			return f.OrganizationID == "64dbf7c7c12e2c3a4b5a6c7d" &&
				len(f.Tags) == 2 &&
				f.Tags[0] == "go" &&
				f.Tags[1] == "api" &&
				f.Page == 2 &&
				f.Limit == 20 &&
				f.SortBy == "created_At" &&
				f.SortOrder == "ASC"
		})).
		Return(expected, total, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/notices?organizationId=64dbf7c7c12e2c3a4b5a6c7d&tags=go,api&page=2&limit=20&sortBy=created_At&sortOrder=ASC", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	var resp struct {
		Data  []NoticeDTO `json:"data"`
		Total int64       `json:"total"`
		Page  int64       `json:"page"`
		Limit int64       `json:"limit"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	s.Len(resp.Data, 1)
	s.Equal("N1", resp.Data[0].Title)
	s.Equal(int64(1), resp.Total)
	s.Equal(int64(2), resp.Page)
	s.Equal(int64(20), resp.Limit)

	s.mockUC.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestUpdateNotice() {
	id := "64dbf7c7c12e2c3a4b5a6c7d"
	body := NoticeDTO{
		Title:   "Updated",
		Content: "New",
		Tags:    []string{"a", "b"},
	}
	jsonBody, _ := json.Marshal(body)

	s.mockUC.On("UpdateNotice", mock.Anything, id, mock.AnythingOfType("*domain.Notice")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/notices/"+id, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.mockUC.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestDeleteNotice() {
	id := "64dbf7c7c12e2c3a4b5a6c7d"
	s.mockUC.On("DeleteNotice", mock.Anything, id).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/notices/"+id, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.mockUC.AssertExpectations(s.T())
}
