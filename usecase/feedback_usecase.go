package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"time"
)

type feedbackUsecase struct {
	feedbackRepo   domain.IFeedbackRepository
	procedureRepo  domain.IProcedureRepository
	contextTimeout time.Duration
}

func NewFeedbackUsecase(fr domain.IFeedbackRepository, pr domain.IProcedureRepository, contextTimeout time.Duration) domain.IFeedbackUsecase {
	return &feedbackUsecase{
		feedbackRepo:   fr,
		procedureRepo:  pr,
		contextTimeout: contextTimeout,
	}
}

func (fu *feedbackUsecase) SubmitFeedback(c context.Context, feedback *domain.Feedback) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	_, err := fu.procedureRepo.GetByID(ctx, feedback.ProcedureID)
	if err != nil {
		return err
	}
	return fu.feedbackRepo.SubmitFeedback(ctx, feedback)
}

func (fu *feedbackUsecase) GetAllFeedbacksForProcedure(c context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	return fu.feedbackRepo.GetAllFeedbacksForProcedure(ctx, procedureID, filter)
}

func (fu *feedbackUsecase) UpdateFeedbackStatus(ctx context.Context, feedbackID, userID string, status domain.FeedbackStatus, adminResponse *string) error {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	feedback, err := fu.feedbackRepo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		return err
	}

	pid := feedback.ProcedureID
	procedure, err := fu.procedureRepo.GetByID(ctx, pid)
	if err != nil {
		return err
	}

	if procedure.OrganizationID != userID {
		return domain.ErrPermissionDenied
	}

	if (status == domain.DeclinedFeedback || status == domain.ResolvedFeedback) && adminResponse == nil {
		return fmt.Errorf("admin response is required when declining or resolving feedback")
	}

	feedback.Status = status
	feedback.AdminResponse = *adminResponse
	feedback.UpdatedAT = time.Now()

	return fu.feedbackRepo.UpdateFeedbackStatus(ctx, feedbackID, feedback)
}

func (fu *feedbackUsecase) GetAllFeedbacks(c context.Context, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.feedbackRepo.GetAllFeedbacks(ctx, filter)
}
