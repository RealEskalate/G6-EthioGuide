package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

type ProcedureUsecase struct {
	procedureRepo  domain.IProcedureRepository
	contextTimeout time.Duration
}

func NewProcedureUsecase(pr domain.IProcedureRepository, timeout time.Duration) domain.IProcedureUsecase {
	return &ProcedureUsecase{
		procedureRepo:  pr,
		contextTimeout: timeout,
	}
}

func (pu *ProcedureUsecase) CreateProcedure(c context.Context, procedure *domain.Procedure) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	return pu.procedureRepo.Create(ctx, procedure)
}

func (pu *ProcedureUsecase) GetProcedureByID(ctx context.Context, id string) (*domain.Procedure, error) {
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

	// 2. Authorization Check:
	// Get userRole and organizationID from context (set by Gin middleware)
	userRole, _ := ctx.Value("userRole").(domain.Role)
	userOrgID, _ := ctx.Value("userID").(string)

	switch userRole {
	case domain.RoleAdmin:
	case domain.RoleOrg:
		if procedureToUpdate.OrganizationID != userOrgID {
			return domain.ErrPermissionDenied
		}
	default:
		return domain.ErrPermissionDenied
	}

	// 3. Update the procedure.
	err = pu.procedureRepo.Update(ctx, id, procedure)
	if err != nil {
		return err
	}

	return nil
}

func (pu *ProcedureUsecase) DeleteProcedure(ctx context.Context, id string) error {
	// 1. Fetch the existing procedure.
	procedureToDelete, err := pu.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Authorization Check:
	userRole, _ := ctx.Value("userRole").(domain.Role)
	userOrgID, _ := ctx.Value("userID").(string)

	switch userRole {
	case domain.RoleAdmin:
	case domain.RoleOrg:
		if procedureToDelete.OrganizationID != userOrgID {
			return domain.ErrPermissionDenied
		}
	default:
		return domain.ErrPermissionDenied
	}

	// 3. Delete the procedure.
	err = pu.procedureRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
