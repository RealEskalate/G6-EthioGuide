package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	Description  string                  `bson:"description,omitempty"`
	Location     string                  `bson:"location,omitempty"`
	Type         domain.OrganizationType `bson:"type,omitempty"` // e.g., "Corporation", "Non-Profit"
	ContactInfo  ContactInfoModel        `bson:"contact_info,omitempty"`
	PhoneNumbers []string                `bson:"phone_numbers,omitempty"`
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

func fromUserDetailDomain(ud *domain.UserDetail) *UserDetailModel {
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

func fromOrgDetailDomain(od *domain.OrganizationDetail) *OrganizationDetailModel {
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

func fromAccountDomain(a domain.Account) AccountModel {
	var objectID primitive.ObjectID
	if id, err := primitive.ObjectIDFromHex(a.ID); err == nil {
		objectID = id
	}

	return AccountModel{
		ID:                 objectID,
		Name:               a.Name,
		Email:              a.Email,
		PasswordHash:       a.PasswordHash,
		Role:               a.Role,
		ProfilePicURL:      a.ProfilePicURL,
		CreatedAt:          a.CreatedAt,
		UserDetail:         fromUserDetailDomain(a.UserDetail),
		OrganizationDetail: fromOrgDetailDomain(a.OrganizationDetail),
	}
}

func toUserDetailDomain(ud UserDetailModel) *domain.UserDetail {
	return &domain.UserDetail{
		Username:         ud.Username,
		SubscriptionPlan: ud.SubscriptionPlan,
		IsBanned:         ud.IsBanned,
		IsVerified:       ud.IsVerified,
	}
}

func toOrgDetailDomain(od OrganizationDetailModel) *domain.OrganizationDetail {
	return &domain.OrganizationDetail{
		Description:  od.Description,
		Location:     od.Location,
		Type:         domain.OrganizationType(od.Type),
		ContactInfo:  domain.ContactInfo(od.ContactInfo),
		PhoneNumbers: od.PhoneNumbers,
	}
}

func toAccountDomain(a AccountModel) *domain.Account {
	return &domain.Account{
		ID:                 a.ID.Hex(),
		Name:               a.Name,
		Email:              a.Email,
		ProfilePicURL:      a.ProfilePicURL,
		Role:               a.Role,
		CreatedAt:          a.CreatedAt,
		UserDetail:         toUserDetailDomain(*a.UserDetail),
		OrganizationDetail: toOrgDetailDomain(*a.OrganizationDetail),
	}
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new user repository instance.
func NewMongoUserRepository(db *mongo.Database, collectionName string) domain.IUserRepository {
	return &MongoUserRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, account *domain.Account) error {
	mongoModel := fromAccountDomain(*account)
	now := time.Now()
	mongoModel.CreatedAt = now
	mongoModel.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, mongoModel)
	if err != nil {
		return err
	}

	// Update the domain object with the generated ID
	account.ID = mongoModel.ID.Hex()
	return nil
}

func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var mongoModel AccountModel
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not an application error, just means no user found
		}
		return nil, err
	}
	return toAccountDomain(mongoModel), nil
}

// GetByUsername fetches a single user by their username.
func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var mongoModel AccountModel
	err := r.collection.FindOne(ctx, bson.M{"user_detail.username": username}).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	// fmt.Println("r", mongoModel.IsActive)
	return toAccountDomain(mongoModel), nil
}

func (dto *AccountModel) ToAccountDomain() *domain.Account {
	return &domain.Account{
		ID:            dto.ID.Hex(),
		Name:          dto.Name,
		Email:         dto.Email,
		PasswordHash:  dto.PasswordHash,
		ProfilePicURL: dto.ProfilePicURL,
		Role:          dto.Role,
		CreatedAt:     dto.CreatedAt,
		UserDetail: &domain.UserDetail{
			SubscriptionPlan: dto.UserDetail.SubscriptionPlan,
			IsBanned:         dto.UserDetail.IsBanned,
			IsVerified:       dto.UserDetail.IsVerified,
		},
		OrganizationDetail: &domain.OrganizationDetail{
			Description: dto.OrganizationDetail.Description,
			Location:    dto.OrganizationDetail.Location,
			Type:        dto.OrganizationDetail.Type,
			ContactInfo: domain.ContactInfo{
				Socials: dto.OrganizationDetail.ContactInfo.Socials,
				Website: dto.OrganizationDetail.ContactInfo.Website,
			},
			PhoneNumbers: dto.OrganizationDetail.PhoneNumbers,
		},
	}
}

func (r *MongoUserRepository) GetByPhone(ctx context.Context, phoneNumber string) (*domain.Account, error) {
	var account AccountModel
	err := r.collection.FindOne(ctx, bson.M{"phone_numbers": phoneNumber}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return account.ToAccountDomain(), err
}
