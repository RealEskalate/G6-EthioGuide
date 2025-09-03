package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockNoticeRepository struct{ mock.Mock }

func (m *MockNoticeRepository) Create(ctx context.Context, notice *domain.Notice) error {
	return m.Called(ctx, notice).Error(0)
}
func (m *MockNoticeRepository) GetByFilter(ctx context.Context, filter *domain.NoticeFilter) ([]*domain.Notice, error) {
	args := m.Called(ctx, filter)
	var res []*domain.Notice
	if args.Get(0) != nil {
		res = args.Get(0).([]*domain.Notice)
	}
	return res, args.Error(1)
}
func (m *MockNoticeRepository) CountByFilter(ctx context.Context, filter *domain.NoticeFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockNoticeRepository) Update(ctx context.Context, id string, notice *domain.Notice) error {
	return m.Called(ctx, id, notice).Error(0)
}
func (m *MockNoticeRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type NoticeUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockNoticeRepository
	uc       *NoticeUsecase
}

func (s *NoticeUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockNoticeRepository)
	s.uc = NewNoticeUsecase(s.mockRepo)
}

func TestNoticeUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(NoticeUsecaseTestSuite))
}

func (s *NoticeUsecaseTestSuite) TestCreateNotice() {
	notice := &domain.Notice{Title: "N1"}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, notice).Return(nil).Once()

		err := s.uc.CreateNotice(context.Background(), notice)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("RepoError", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, notice).Return(errors.New("db")).Once()

		err := s.uc.CreateNotice(context.Background(), notice)
		s.Error(err)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *NoticeUsecaseTestSuite) TestGetNoticesByFilter() {
	s.Run("Enforces defaults and caps limit", func() {
		filter := &domain.NoticeFilter{Page: 0, Limit: 1000}

		s.mockRepo.On("GetByFilter", mock.Anything, mock.MatchedBy(func(f *domain.NoticeFilter) bool {
			return f.Page == 1 && f.Limit == 100
		})).Return([]*domain.Notice{}, nil).Once()
		s.mockRepo.On("CountByFilter", mock.Anything, mock.Anything).Return(int64(0), nil).Once()

		_, total, err := s.uc.GetNoticesByFilter(context.Background(), filter)
		s.NoError(err)
		s.Equal(int64(0), total)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *NoticeUsecaseTestSuite) TestUpdateNotice() {
	notice := &domain.Notice{Title: "Updated"}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Update", mock.Anything, "id123", notice).Return(nil).Once()

		err := s.uc.UpdateNotice(context.Background(), "id123", notice)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("RepoError", func() {
		s.SetupTest()
		s.mockRepo.On("Update", mock.Anything, "id123", notice).Return(errors.New("db")).Once()

		err := s.uc.UpdateNotice(context.Background(), "id123", notice)
		s.Error(err)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *NoticeUsecaseTestSuite) TestDeleteNotice() {
	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Delete", mock.Anything, "id123").Return(nil).Once()

		err := s.uc.DeleteNotice(context.Background(), "id123")
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("RepoError", func() {
		s.SetupTest()
		s.mockRepo.On("Delete", mock.Anything, "id123").Return(errors.New("db")).Once()

		err := s.uc.DeleteNotice(context.Background(), "id123")
		s.Error(err)
		s.mockRepo.AssertExpectations(s.T())
	})
}
