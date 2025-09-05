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

	checklist, err := cuc.checklistrepo.FindCheck(ctx, checklistID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checklist: %w", err)
	}

	if errToggle := cuc.checklistrepo.ToggleCheck(ctx, checklistID); errToggle != nil {
		return nil, fmt.Errorf("failed to update checklist: %w", errToggle)
	}

	filterGeneral := interface{}(map[string]interface{}{
		"user_procedure_id": checklist.UserProcedureID,
	})

	filterChecked := interface{}(map[string]interface{}{
		"user_procedure_id": checklist.UserProcedureID,
		"is_checked":        true,
	})

	countDoc, errDoc := cuc.checklistrepo.CountDocumentsChecklist(ctx, filterGeneral)
	if errDoc != nil {
		return nil, fmt.Errorf("failed to update fields: %w", errDoc)
	}

	countChecked, errChecked := cuc.checklistrepo.CountDocumentsChecklist(ctx, filterChecked)
	if errChecked != nil {
		return nil, fmt.Errorf("failed to update fields: %w", errChecked)
	}

	updatefields := make(map[string]interface{})
	percent := int((float64(countChecked) / float64(countDoc)) * 100)
	if percent == 0 {
		updatefields["status"] = "Not Started"
	} else if percent > 0 && percent < 100 {
		updatefields["status"] = "In Progress"
	} else {
		updatefields["status"] = "Completed"
	}

	filterUserProcedure := interface{}(map[string]interface{}{
		"user_procedure_id": checklist.UserProcedureID,
	})

	if err := cuc.checklistrepo.UpdateUserProcedure(ctx, filterUserProcedure, updatefields); err != nil {
		return nil, fmt.Errorf("failed to update fields: %w", errChecked)
	}

	checklist.IsChecked = !checklist.IsChecked
	return checklist, nil
}
