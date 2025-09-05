package domain

import "time"

type UserProcedure struct {
	ID          string
	UserID      string
	ProcedureID string
	Percent     int
	Status      string // e.g., "Not Started", "In Progress", "Completed"
	UpdatedAt   time.Time
}

