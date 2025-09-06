package domain

import "time"

type AIProcedure struct {
	Id   string
	Name string
}

type AIChat struct {
	ID                string
	UserID            string
	Source            string
	Request           string
	Response          string
	Timestamp         time.Time
	RelatedProcedures []*AIProcedure
}
