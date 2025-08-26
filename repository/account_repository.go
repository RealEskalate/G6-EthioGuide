package repository

import (
	"EthioGuide/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactInfoModel struct {
	Socials map[string]string `bson:"socials,omitempty"`
	Website string            `bson:"website,omitempty"`
}

type UserDetailModel struct {
	Username         string              `bson:"username"`
	SubscriptionPlan domain.Subscription `bson:"subscription_plan,omitempty"`
	IsBanned         bool                `bson:"is_banned"`
	IsVerified       bool                `bson:"is_verified"`
}

type OrganizationDetailModel struct {
	Description  string           `bson:"description,omitempty"`
	Location     string           `bson:"location,omitempty"`
	Type         string           `bson:"type,omitempty"` // e.g., "Corporation", "Non-Profit"
	ContactInfo  ContactInfoModel `bson:"contact_info,omitempty"`
	PhoneNumbers []string         `bson:"phone_numbers,omitempty"`
}

type AccountModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Email         string             `bson:"email"`
	PasswordHash  string             `bson:"password_hash"`
	ProfilePicURL string             `bson:"profile_pic_url,omitempty"`
	Role          domain.Role        `bson:"role"`
	CreatedAt     time.Time          `bson:"created_at"`

	UserDetail         *UserDetailModel         `bson:"user_detail,omitempty"`
	OrganizationDetail *OrganizationDetailModel `bson:"organization_detail,omitempty"`
}
