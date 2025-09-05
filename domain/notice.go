package domain

import "time"

type Notice struct {
	ID             string
	OrganizationID string
	Title          string
	Content        string
	Tags           []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NoticeFilter struct {
	// filter
	OrganizationID string
	Tags           []string

	// sort
	SortBy    string
	SortOrder SortOrder

	// Pagination
	Page  int64
	Limit int64
}
