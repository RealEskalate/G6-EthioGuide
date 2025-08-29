package domain

type Category struct {
	ID             string
	OrganizationID string
	ParentID       string
	Title          string
}