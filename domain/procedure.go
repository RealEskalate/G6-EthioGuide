package domain

import "time"

type Content struct {
	Prerequisites []string
	Steps         []string
	Result        []string
}

type Fees struct {
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
	GroupID        string
	OrganizationID string
	Name           string
	Content        Content
	Fees           Fees
	ProcessingTime ProcessingTime
	CreatedAt      time.Time
	UpdatedAt      time.Time
}