package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

// SearchUsecase implements the domain.ISearchUsecase interface.
type SearchUsecase struct {
	searchRepo     domain.ISearchRepository
	contextTimeout time.Duration
}

// NewSearchUsecase is the constructor for SearchUsecase.
func NewSearchUsecase(sr domain.ISearchRepository, timeout time.Duration) *SearchUsecase {
	return &SearchUsecase{
		searchRepo:     sr,
		contextTimeout: timeout,
	}
}

// Search orchestrates the search operation by validating input
// and calling the repository.
func (uc *SearchUsecase) Search(c context.Context, filters domain.SearchFilterRequest) (*domain.SearchResult, error) {
	// Create a context with a timeout to prevent long-running queries.
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// --- Validation and Defaults ---
	// This is a key responsibility of the use case layer.
	if filters.Limit <= 0 || filters.Limit > 50 {
		filters.Limit = 20 // Set a default and max limit
	}
	if filters.Page <= 0 {
		filters.Page = 1 // Default to the first page
	}

	// --- Call the Repository ---
	// The use case delegates the actual data fetching to the repository.
	searchResults, err := uc.searchRepo.Search(ctx, filters)
	if err != nil {
		// For now, we just pass the error up. In the future, you could
		// add specific error handling or logging here.
		return nil, err
	}

	return searchResults, nil
}