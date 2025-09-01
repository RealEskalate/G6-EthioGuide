package domain

type Category struct {
	ID             string
	OrganizationID string
	ParentID       string
	Title          string
}

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type CategorySearchAndFilter struct {
	Title string

	ParentID       string
	OrganizationID string

	Page  int64
	Limit int64

	SortBy    string
	SortOrder SortOrder
}