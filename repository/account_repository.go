package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactInfo struct {
	Socials map[string]string `bson:"socials,omitempty"`
	Website string            `bson:"website,omitempty"`
}

type UserDetailModel struct {
	SubscriptionPlan domain.Subscription `bson:"subscription_plan,omitempty"`
	IsBanned         bool                `bson:"is_banned"`
	IsVerified       bool                `bson:"is_verified"`
}

type OrganizationDetailModel struct {
	Description  string      `bson:"description,omitempty"`
	Location     string      `bson:"location,omitempty"`
	Type         domain.OrganizationType      `bson:"type,omitempty"` 
	ContactInfo  ContactInfo `bson:"contact_info,omitempty"`
	PhoneNumbers []string    `bson:"phone_numbers,omitempty"`
}

type AccountModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Email         string             `bson:"email"`
	PasswordHash  string             `bson:"password_hash"`
	ProfilePicURL string             `bson:"profile_pic_url,omitempty"`
	Role          domain.Role        `bson:"role"`
	CreatedAt     time.Time          `bson:"created_at"`

	UserDetail         *UserDetailModel        `bson:"user_detail,omitempty"`
	OrganizationDetail *OrganizationDetailModel `bson:"organization_detail,omitempty"`
}

func (dto *AccountModel) ToAccountDomain() *domain.Account {
	return &domain.Account{
		ID:            dto.ID.Hex(),
		Name:          dto.Name,
		Email:         dto.Email,
		PasswordHash:  dto.PasswordHash,
		ProfilePicURL: dto.ProfilePicURL,
		Role:         dto.Role,
		CreatedAt:    dto.CreatedAt,
		UserDetail:  &domain.UserDetail{
			SubscriptionPlan: dto.UserDetail.SubscriptionPlan,
			IsBanned: dto.UserDetail.IsBanned,
			IsVerified: dto.UserDetail.IsVerified,
		},
		OrganizationDetail: &domain.OrganizationDetail{
			Description:  dto.OrganizationDetail.Description,
			Location:     dto.OrganizationDetail.Location,
			Type:        dto.OrganizationDetail.Type,
			ContactInfo: domain.ContactInfo{
				Socials: dto.OrganizationDetail.ContactInfo.Socials,
				Website: dto.OrganizationDetail.ContactInfo.Website,
			},
			PhoneNumbers: dto.OrganizationDetail.PhoneNumbers,
		},
	}
}

type AccountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepository(db *mongo.Database) *AccountRepository {
	collection := db.Collection("accounts")
	return &AccountRepository{
		collection: collection,
	}
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var account AccountModel
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return account.ToAccountDomain(), err
}

func (r *AccountRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Account, error) {
	var account AccountModel
	err := r.collection.FindOne(ctx, bson.M{"phone_numbers": phoneNumber}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return account.ToAccountDomain(), err
}

func (r *AccountRepository) FindByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var account AccountModel
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return account.ToAccountDomain(), err
}