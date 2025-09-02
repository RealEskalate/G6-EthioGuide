package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

type procedureUsecase struct {
	procedureRepo domain.IProcedureRepository
	contextTimeout time.Duration
}

func NewProcedureUsecase(pr domain.IProcedureRepository, timeout time.Duration) domain.IProcedureUsecase {
	return &procedureUsecase{
		procedureRepo: pr,
		contextTimeout: timeout,
	}
}

func (pu *procedureUsecase) CreateProcedure(c context.Context, procedure *domain.Procedure) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	return pu.procedureRepo.Create(ctx, procedure)
}