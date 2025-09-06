package domain

import (
	"fmt"
	"sort"
	"strings"
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
	Embedding      []float64
	NoticeIDs      []string
}

func (p Procedure) ToString() string {
	var sb strings.Builder

	// 1. Name: The primary identifier of the procedure's purpose.
	// We use a descriptive prefix to give the model context.
	fmt.Fprintf(&sb, "Procedure Name: %s\n\n", p.Name)

	// 2. Content: The core of what the procedure entails.
	// Prerequisites
	if len(p.Content.Prerequisites) > 0 {
		sb.WriteString("Prerequisites:\n")
		for _, prereq := range p.Content.Prerequisites {
			fmt.Fprintf(&sb, "- %s\n", prereq)
		}
		sb.WriteString("\n")
	}

	// Steps (sorted for logical flow)
	if len(p.Content.Steps) > 0 {
		sb.WriteString("Steps to complete:\n")

		// Sort keys to ensure a consistent, logical order.
		keys := make([]int, 0, len(p.Content.Steps))
		for k := range p.Content.Steps {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			fmt.Fprintf(&sb, "%d. %s\n", k, p.Content.Steps[k])
		}
		sb.WriteString("\n")
	}

	// Result
	if p.Content.Result != "" {
		fmt.Fprintf(&sb, "Expected Result: %s\n\n", p.Content.Result)
	}

	// 3. Key Attributes: Important properties that define the procedure.
	// Fees
	// Using the label provides crucial context for the cost.
	fmt.Fprintf(&sb, "Fee: %s - %.2f %s\n", p.Fees.Label, p.Fees.Amount, p.Fees.Currency)

	// Processing Time
	sb.WriteString("Processing Time: ")
	if p.ProcessingTime.MinDays == p.ProcessingTime.MaxDays {
		fmt.Fprintf(&sb, "%d days\n", p.ProcessingTime.MinDays)
	} else {
		fmt.Fprintf(&sb, "%d to %d days\n", p.ProcessingTime.MinDays, p.ProcessingTime.MaxDays)
	}

	// Return the final, cleaned-up string.
	// TrimSpace removes any leading/trailing newlines.
	return strings.TrimSpace(sb.String())
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
	Name *string // search in name

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
