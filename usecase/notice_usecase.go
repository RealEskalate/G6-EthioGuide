package usecase

import (
	"EthioGuide/domain"
	"context"
)

type NoticeUsecase struct {
	noticeRepo domain.INoticeRepository
}

func NewNoticeUsecase(noticeRepo domain.INoticeRepository) *NoticeUsecase {
	return &NoticeUsecase{
		noticeRepo: noticeRepo,
	}
}

func (nu *NoticeUsecase) CreateNotice(ctx context.Context, notice *domain.Notice) error {
	err := nu.noticeRepo.Create(ctx, notice)
	if err != nil {
		return err
	}
	return nil
}

func (nu *NoticeUsecase) GetNoticesByFilter(ctx context.Context, filter *domain.NoticeFilter) ([]*domain.Notice, int64, error) {
       // Enforce pagination defaults/bounds
       if filter.Limit <= 0 {
	       filter.Limit = 10
       }
       if filter.Limit > 100 {
	       filter.Limit = 100
       }
       if filter.Page <= 0 {
	       filter.Page = 1
       }

       // Only created_At sorting is supported; ensure expected defaults
       if filter.SortBy == "" {
	       filter.SortBy = "created_At"
       }
       if filter.SortOrder != domain.SortOrderASC && filter.SortOrder != domain.SortOrderDESC {
	       filter.SortOrder = domain.SortOrderDESC
       }

       notices, err := nu.noticeRepo.GetByFilter(ctx, filter)
       if err != nil {
	       return nil, 0, err
       }
       total, err := nu.noticeRepo.CountByFilter(ctx, filter)
       if err != nil {
	       return nil, 0, err
       }
       return notices, total, nil
}

func (nu *NoticeUsecase) UpdateNotice(ctx context.Context, id string, notice *domain.Notice) error {
	err := nu.noticeRepo.Update(ctx, id, notice)
	if err != nil {
		return err
	}
	return nil
}


func (un *NoticeUsecase) DeleteNotice(ctx context.Context, id string) error {
	err := un.noticeRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}