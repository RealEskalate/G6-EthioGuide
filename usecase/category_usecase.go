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

func (cc *categoryUsecase) GetCategories(c context.Context, options *domain.CategorySearchAndFilter) ([]*domain.Category, int64, error) {
	ctx, cancel := context.WithTimeout(c, cc.contextTimeout)
	defer cancel()

	if options.Limit <= 0 {
		options.Limit = 10
	}
	if options.Limit > 100 {
		options.Limit = 100 // Enforce a max limit
	}
	if options.Page <= 0 {
		options.Page = 1
	}

	return cc.categoryRepo.GetCategories(ctx, options)
}