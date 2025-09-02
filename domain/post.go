package domain

import "time"

type Post struct {
	UserID    string
	Title    string
	Content   string
	Procedures     []string
	Tags 			[]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

