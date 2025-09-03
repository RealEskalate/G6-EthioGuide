package domain

import "time"

type FeedbackType string

const (
	InaccuracyFeedback FeedbackType = "inaccuracy"
	OutdatedFeedback   FeedbackType = "outdated"
	ThanksFeedback     FeedbackType = "thanks"
	MissingFeedback    FeedbackType = "missing"
)

type FeedbackStatus string

const (
	NewFeedback        FeedbackStatus = "new"
	InProgressFeedback FeedbackStatus = "in_progress"
	ResolvedFeedback   FeedbackStatus = "resolved"
	DeclinedFeedback   FeedbackStatus = "declined"
)

type Feedback struct {
	ID            string
	UserID        string
	ProcedureID   string
	Content       string
	LikeCount     int
	DislikeCount  int
	ViewCount     int
	Type          FeedbackType
	Status        FeedbackStatus
	AdminResponse string
	Tags          []string
	CreatedAT     time.Time
	UpdatedAT     time.Time
}

type FeedbackFilter struct {
	Page int64
	Limit int64
	Status *string
}