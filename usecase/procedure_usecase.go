package usecase

import (
	"EthioGuide/domain"
	"context"
)

type ProcedureUsecase struct {
	procedureRepo domain.IProcedureRepository
}

func NewProcedureUsecase(procedureRepo domain.IProcedureRepository) *ProcedureUsecase {
	return &ProcedureUsecase{
		procedureRepo: procedureRepo,
	}
}

func (pu *ProcedureUsecase) GetProcedureByID(ctx context.Context, id string) (*domain.Procedure, error){
	procedure, err := pu.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return procedure, nil
}

func (pu *ProcedureUsecase) UpdateProcedure(ctx context.Context, id string, procedure *domain.Procedure) error {
	// 1. Fetch the existing procedure.
	procedureToUpdate, err := pu.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return err 
	}

	// 2. Authorization Check: only the organization can update their post.
	if procedureToUpdate.OrganizationID != procedure.OrganizationID {
		return domain.ErrPermissionDenied
	}

	// 3. Update the procedure.
	err = pu.procedureRepo.Update(ctx, id, procedure)
	if err != nil {
		return err
	}

	return nil
}