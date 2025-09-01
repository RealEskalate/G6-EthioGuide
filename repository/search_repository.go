package repository

import "context"

type SearchRepository struct {
	ProcedureRepo *ProcedureRepository
	AccountRepo   *AccountRepository
}

func NewSearchRepository(prodrepo *ProcedureRepository, accrepo *AccountRepository) *SearchRepository {
	return &SearchRepository{
		ProcedureRepo: prodrepo,
		AccountRepo:   accrepo,
	}
}

func (r *SearchRepository) Search(ctx context.Context, filter domain)