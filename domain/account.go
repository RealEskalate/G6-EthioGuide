package domain

type Role string
type Subscription string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleOrg   Role = "org"

	SubscriptionNone Subscription = "none"
	SubscriptionPro  Subscription = "pro"
)

// IsValid checks if the role is one of the predefined valid roles.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleOrg:
		return true
	}
	return false
}
