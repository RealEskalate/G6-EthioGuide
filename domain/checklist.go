package domain

type Checklist struct {
	ID              string
	UserProcedureID string
	Type            string
	Content         string
	IsChecked       bool
}
