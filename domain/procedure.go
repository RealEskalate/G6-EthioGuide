package domain

import (
	"time"
)

type ProcedureContent struct {
	Prerequisites []string
	Steps         map[int]string
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


// ====== Search & Filter Options ======
type GlobalLogic string
// type SortOrder string

const (
	GlobalLogicOR  GlobalLogic = "OR"
	GlobalLogicAND GlobalLogic = "AND"
)

type ProcedureSearchFilterOptions struct {
	// Search
	Name     *string // search in name

	// Filters
	OrganizationID *string
	GroupID        *string

	// Fee range
	MinFee *float64
	MaxFee *float64

	// Processing time range
	MinProcessingDays *int
	MaxProcessingDays *int

	// Date range
	StartDate *time.Time
	EndDate   *time.Time

	// Logic
	GlobalLogic GlobalLogic

	// Pagination
	Page  int64
	Limit int64

	// Sorting
	SortBy    string
	SortOrder SortOrder
}