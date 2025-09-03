package domain

import (
	"time"
)

type ProcedureContent struct {
	Prerequisites []string
	Steps         []string
	Result        string
}

type ProcedureFee struct {
	Label    string
	Currency string
	Amount   float64
}

type ProcessingTime struct {
	MinDays int
	MaxDays int
}

type Procedure struct {
	ID             string
	GroupID        *string
	OrganizationID string
	Name           string
	Content        ProcedureContent
	Fees           ProcedureFee
	ProcessingTime ProcessingTime
	CreatedAt      time.Time
	NoticeIDs      []string
}