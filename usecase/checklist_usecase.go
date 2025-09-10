package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"time"
)

type ChecklistUsecase struct {
	checklistrepo domain.IChecklistRepository
}

func NewChecklistUsecase(checkrepo domain.IChecklistRepository) *ChecklistUsecase {
	return &ChecklistUsecase{
		checklistrepo: checkrepo,
	}
}

func (cuc *ChecklistUsecase) CreateChecklist(ctx context.Context, userid, procedureID string) (*domain.UserProcedure, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userid == "" || procedureID == "" {
		return nil, domain.ErrInvalidID
	}

	userProcedure, err := cuc.checklistrepo.CreateChecklist(ctx, userid, procedureID)
	if err != nil {
		return nil, fmt.Errorf("failed to create checklist in usecase layer for user %s: %w", userid, err)
	}

	return userProcedure, nil
}

func (cuc *ChecklistUsecase) GetProcedures(ctx context.Context, userid string) ([]*domain.UserProcedure, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userid == "" {
		return nil, domain.ErrInvalidID
	}

	procedures, err := cuc.checklistrepo.GetProcedures(ctx, userid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user procedures: %w", err)
	}

	return procedures, nil

}

func (cuc *ChecklistUsecase) GetChecklistByUserProcedureID(ctx context.Context, userprocedureID string) ([]*domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userprocedureID == "" {
		return nil, domain.ErrInvalidID
	}

	checklists, err := cuc.checklistrepo.GetChecklistByUserProcedureID(ctx, userprocedureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checklists: %w", err)
	}

	return checklists, nil
}

func (cuc *ChecklistUsecase) UpdateChecklist(ctx context.Context, checklistID string) (*domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if checklistID == "" {
		return nil, domain.ErrInvalidID
	}

	updatedChecklist, err := cuc.checklistrepo.ToggleCheckAndUpdateStatus(ctx, checklistID)
	if err != nil {
		return nil, err
	}

	return updatedChecklist, nil
}
