package domain

import "time"

type Role string
type Subscription string
type OrganizationType string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleOrg   Role = "org"

	SubscriptionNone Subscription = "none"
	SubscriptionPro  Subscription = "pro"

	OrgTypeGov     OrganizationType = "gov"
	OrgTypePrivate OrganizationType = "private"
)

// IsValid checks if the role is one of the predefined valid roles.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleOrg:
		return true
	}
	return false
}

func (s Subscription) IsValid() bool {
	switch s {
	case SubscriptionPro, SubscriptionNone:
		return true
	}
	return false
}

func (t OrganizationType) IsValid() bool {
	switch t {
	case OrgTypeGov, OrgTypePrivate:
		return true
	}
	return false
}

type Account struct {
	ID            string
	Name          string
	Email         string
	PasswordHash  string
	ProfilePicURL string
	Role          Role
	CreatedAt     time.Time

	UserDetail         *UserDetail
	OrganizationDetail *OrganizationDetail
}

type UserDetail struct {
	Username         string
	SubscriptionPlan Subscription
	IsBanned         bool
	IsVerified       bool
}

type OrganizationDetail struct {
	Description  string
	Location     string
	Type         OrganizationType
	ContactInfo  ContactInfo
	PhoneNumbers []string
}

type ContactInfo struct {
	Socials map[string]string
	Website string
}

type GetOrgsFilter struct {
	Type     string
	Query    string
	Page     int64
	PageSize int64
}
