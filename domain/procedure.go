package domain

import "time"

type Procedure struct {
	ID             int
	GroupID        *int
	OrganizationID int
	Name           string
	Content        map[string]interface{}
	Fees           map[string]interface{}
	ProcessingTime map[string]interface{}
	CreatedAt      time.Time
	NoticeIDs      []string
}
