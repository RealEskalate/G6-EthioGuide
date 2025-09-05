package domain

import "time"

type Post struct {
	ID         string
	UserID     string
	Title      string
	Content    string
	Procedures []string
	Tags       []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// type GlobalLogic string
// type ActionType string

type PostFilters struct {
	Title *string

	// List of tags
	ProcedureID []string
	Tags        []string

	StartDate *time.Time
	EndDate   *time.Time

	Page  int64
	Limit int64

	SortBy string
	// ASC or DESC
	SortOrder SortOrder
}

func (p PostFilters) Find() any {
	panic("unimplemented")
}
