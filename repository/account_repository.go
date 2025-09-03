package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// --- Database Models ---

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
	Description  string                  `bson:"description,omitempty"`
	Location     string                  `bson:"location,omitempty"`
	Type         domain.OrganizationType `bson:"type,omitempty"`
	ContactInfo  ContactInfoModel        `bson:"contact_info,omitempty"`
	PhoneNumbers []string                `bson:"phone_numbers,omitempty"`
}

type AccountModel struct {
	ID                 primitive.ObjectID       `bson:"_id,omitempty"`
	Name               string                   `bson:"name,omitempty"`
	Email              string                   `bson:"email"`
	PasswordHash       string                   `bson:"password_hash"`
	ProfilePicURL      string                   `bson:"profile_pic_url,omitempty"`
	Role               domain.Role              `bson:"role"`
	CreatedAt          time.Time                `bson:"created_at"`
	UserDetail         *UserDetailModel         `bson:"user_detail,omitempty"`
	OrganizationDetail *OrganizationDetailModel `bson:"organization_detail,omitempty"`
}

// --- Mappers ---

func fromDomainAccount(a *domain.Account) (*AccountModel, error) {
	return &AccountModel{
		Name:               a.Name,
		Email:              a.Email,
		PasswordHash:       a.PasswordHash,
		Role:               a.Role,
		ProfilePicURL:      a.ProfilePicURL,
		CreatedAt:          a.CreatedAt,
		UserDetail:         fromDomainUserDetail(a.UserDetail),
		OrganizationDetail: fromDomainOrgDetail(a.OrganizationDetail),
	}, nil
}

func fromDomainUserDetail(ud *domain.UserDetail) *UserDetailModel {
	if ud == nil {
		return nil
	}
	return &UserDetailModel{
		Username:         ud.Username,
		SubscriptionPlan: ud.SubscriptionPlan,
		IsBanned:         ud.IsBanned,
		IsVerified:       ud.IsVerified,
	}
}

func fromDomainOrgDetail(od *domain.OrganizationDetail) *OrganizationDetailModel {
	if od == nil {
		return nil
	}
	return &OrganizationDetailModel{
		Description:  od.Description,
		Location:     od.Location,
		Type:         od.Type,
		ContactInfo:  ContactInfoModel(od.ContactInfo),
		PhoneNumbers: od.PhoneNumbers,
	}
}

func toDomainAccount(a *AccountModel) *domain.Account {
	domainAccount := &domain.Account{
		ID:            a.ID.Hex(),
		Name:          a.Name,
		Email:         a.Email,
		PasswordHash:  a.PasswordHash,
		ProfilePicURL: a.ProfilePicURL,
		Role:          a.Role,
		CreatedAt:     a.CreatedAt,
	}

	if a.UserDetail != nil {
		domainAccount.UserDetail = toDomainUserDetail(a.UserDetail)
	}

	if a.OrganizationDetail != nil {
		domainAccount.OrganizationDetail = toDomainOrgDetail(a.OrganizationDetail)
	}

	return domainAccount
}

func toDomainUserDetail(ud *UserDetailModel) *domain.UserDetail {
	return &domain.UserDetail{
		Username:         ud.Username,
		SubscriptionPlan: ud.SubscriptionPlan,
		IsBanned:         ud.IsBanned,
		IsVerified:       ud.IsVerified,
	}
}

func toDomainOrgDetail(od *OrganizationDetailModel) *domain.OrganizationDetail {
	return &domain.OrganizationDetail{
		Description:  od.Description,
		Location:     od.Location,
		Type:         od.Type,
		ContactInfo:  domain.ContactInfo(od.ContactInfo),
		PhoneNumbers: od.PhoneNumbers,
	}
}

// --- Repository Implementation ---

type AccountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepository(db *mongo.Database) domain.IAccountRepository {
	return &AccountRepository{
		collection: db.Collection("accounts"),
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	model, err := fromDomainAccount(account)
	if err != nil {
		return fmt.Errorf("failed to map domain account to model: %w", err)
	}

	model.CreatedAt = time.Now()
	model.ID = primitive.NewObjectID()

	_, err = r.collection.InsertOne(ctx, model)
	if err != nil {
		// Here you might check for mongo duplicate key errors and return domain.Err...
		return fmt.Errorf("failed to insert account: %w", err)
	}

	// Back-fill the domain object with the generated ID.
	account.ID = model.ID.Hex()
	return nil
}

func (r *AccountRepository) GetById(ctx context.Context, id string) (*domain.Account, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var model AccountModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomainAccount(&model), nil
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var model AccountModel
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomainAccount(&model), nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var model AccountModel
	err := r.collection.FindOne(ctx, bson.M{"user_detail.username": username}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomainAccount(&model), nil
}

func (r *AccountRepository) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"password_hash": newPassword,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AccountRepository) UpdateProfile(ctx context.Context, account domain.Account) error {
	accountID, err := primitive.ObjectIDFromHex(account.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	updatedAccount, err := fromDomainAccount(&account)
	if err != nil {
		return err
	}
	update := bson.M{"$set": updatedAccount}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": accountID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
