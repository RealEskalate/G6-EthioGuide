package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

type categoryUsecase struct {
	categoryRepo domain.ICategoryRepository
	contextTimeout time.Duration
}

func NewCategoryUsecase(cr domain.ICategoryRepository, timeout time.Duration) domain.ICategoryUsecase {
	return &categoryUsecase{
		categoryRepo: cr,
		contextTimeout: timeout,
	}
}

func (cc *categoryUsecase) CreateCategory(c context.Context, category *domain.Category) error {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	return cc.categoryRepo.Create(ctx, category)
}