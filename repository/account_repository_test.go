package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository" // Dot import for convenience
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.IAccountRepository
	collection *mongo.Collection
}

// SetupSuite connects to the database initialized in TestMain.
func (s *AccountRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewAccountRepository(s.db)
	s.collection = s.db.Collection("accounts")
}

// TearDownSuite drops the database after all tests in this suite are done.
func (s *AccountRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	s.Require().NoError(err)
}

// BeforeTest cleans the collection to ensure test isolation.
func (s *AccountRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err)
}

// TestAccountRepositoryTestSuite is the entry point for the suite.
func TestAccountRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(AccountRepositoryTestSuite))
}

// --- Test Cases ---

func (s *AccountRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	account := &domain.Account{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "some_hash",
		Role:         domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "testuser",
		},
	}

	err := s.repo.Create(ctx, account)

	s.NoError(err)
	s.NotEmpty(account.ID, "Domain account ID should be back-filled after creation")

	// Verify directly in the DB
	var createdModel AccountModel
	objID, _ := primitive.ObjectIDFromHex(account.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
	s.NoError(err)
	s.Equal("test@example.com", createdModel.Email)
	s.Equal("testuser", createdModel.UserDetail.Username)
}

func (s *AccountRepositoryTestSuite) TestGetById() {
	ctx := context.Background()
	// Arrange: Create a user to fetch
	userToCreate := &domain.Account{
		Name:  "Get User",
		Email: "getbyid@example.com",
		Role:  domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "getuser",
		},
	}
	err := s.repo.Create(ctx, userToCreate)
	s.Require().NoError(err, "Setup: failed to create user")

	s.Run("Success", func() {
		foundAccount, err := s.repo.GetById(ctx, userToCreate.ID)
		s.NoError(err)
		s.NotNil(foundAccount)
		s.Equal(userToCreate.ID, foundAccount.ID)
		s.Equal("getbyid@example.com", foundAccount.Email)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		foundAccount, err := s.repo.GetById(ctx, nonExistentID)
		s.Error(err)
		s.Nil(foundAccount)
		s.ErrorIs(err, domain.ErrNotFound, "Should return a domain-specific not found error")
	})

	s.Run("Failure - Invalid ID Format", func() {
		foundAccount, err := s.repo.GetById(ctx, "this-is-not-an-object-id")
		s.Error(err)
		s.Nil(foundAccount)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *AccountRepositoryTestSuite) TestGetByEmail() {
	ctx := context.Background()
	// Arrange
	userToCreate := &domain.Account{
		Email: "findme@example.com",
		Role:  domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "findmeuser",
		},
	}
	err := s.repo.Create(ctx, userToCreate)
	s.Require().NoError(err)

	s.Run("Success", func() {
		foundAccount, err := s.repo.GetByEmail(ctx, "findme@example.com")
		s.NoError(err)
		s.NotNil(foundAccount)
		s.Equal(userToCreate.ID, foundAccount.ID)
	})

	s.Run("Failure - Not Found", func() {
		foundAccount, err := s.repo.GetByEmail(ctx, "dne@example.com")
		s.Error(err)
		s.Nil(foundAccount)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *AccountRepositoryTestSuite) TestGetByUsername() {
	ctx := context.Background()
	// Arrange
	userToCreate := &domain.Account{
		Email: "user@example.com",
		Role:  domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "findmebyusername",
		},
	}
	err := s.repo.Create(ctx, userToCreate)
	s.Require().NoError(err)

	s.Run("Success", func() {
		foundAccount, err := s.repo.GetByUsername(ctx, "findmebyusername")
		s.NoError(err)
		s.NotNil(foundAccount)
		s.Equal(userToCreate.ID, foundAccount.ID)
	})

	s.Run("Failure - Not Found", func() {
		foundAccount, err := s.repo.GetByUsername(ctx, "dne_username")
		s.Error(err)
		s.Nil(foundAccount)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *AccountRepositoryTestSuite) TestMappingLogic() {
	ctx := context.Background()

	s.Run("Correctly maps User account", func() {
		user := &domain.Account{
			Name:       "Mapping User",
			Email:      "mapping.user@example.com",
			Role:       domain.RoleUser,
			UserDetail: &domain.UserDetail{Username: "mappinguser"},
		}
		err := s.repo.Create(ctx, user)
		s.Require().NoError(err)

		foundUser, err := s.repo.GetById(ctx, user.ID)
		s.NoError(err)
		s.NotNil(foundUser.UserDetail, "UserDetail should not be nil")
		s.Nil(foundUser.OrganizationDetail, "OrganizationDetail should be nil for a user")
		s.Equal("mappinguser", foundUser.UserDetail.Username)
	})

	s.Run("Correctly maps Organization account", func() {
		org := &domain.Account{
			Name:               "Mapping Org",
			Email:              "mapping.org@example.com",
			Role:               domain.RoleOrg,
			OrganizationDetail: &domain.OrganizationDetail{Location: "Test Location"},
		}
		err := s.repo.Create(ctx, org)
		s.Require().NoError(err)

		foundOrg, err := s.repo.GetById(ctx, org.ID)
		s.NoError(err)
		s.Nil(foundOrg.UserDetail, "UserDetail should be nil for an org")
		s.NotNil(foundOrg.OrganizationDetail, "OrganizationDetail should not be nil")
		s.Equal("Test Location", foundOrg.OrganizationDetail.Location)
	})
}

func (s *AccountRepositoryTestSuite) TestUpdateProfile() {
	ctx := context.Background()

	// First, create a user to update
	original := &domain.Account{
		Name:         "Original Name",
		Email:        "updateprofile@example.com",
		PasswordHash: "hashedpassword",
		Role:         domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username:         "updateuser",
			SubscriptionPlan: domain.SubscriptionNone,
			IsBanned:         false,
			IsVerified:       true,
		},
	}
	err := s.repo.Create(ctx, original)
	s.Require().NoError(err)
	s.Require().NotEmpty(original.ID)

	s.Run("Success", func() {
		updated := *original
		updated.Name = "Updated Name"
		updated.UserDetail.Username = "updateduser"

		err := s.repo.UpdateProfile(ctx, updated)
		s.NoError(err)

		// Fetch and verify
		fetched, err := s.repo.GetById(ctx, original.ID)
		s.NoError(err)
		s.Equal("Updated Name", fetched.Name)
		s.Equal("updateduser", fetched.UserDetail.Username)
	})

	s.Run("Failure - Invalid ID", func() {
		invalid := *original
		invalid.ID = "badid"
		err := s.repo.UpdateProfile(ctx, invalid)
		s.ErrorIs(err, domain.ErrUserNotFound)
	})

	s.Run("Failure - Not Found", func() {
		notFound := *original
		notFound.ID = "507f1f77bcf86cd799439011" // valid ObjectID, but not in DB
		err := s.repo.UpdateProfile(ctx, notFound)
		s.ErrorIs(err, domain.ErrUserNotFound)
	})
}

func (s *AccountRepositoryTestSuite) TestUpdatePassword() {
	ctx := context.Background()

	// Arrange: Create a user to update password for
	account := &domain.Account{
		Name:         "Password User",
		Email:        "passworduser@example.com",
		PasswordHash: "old_hash",
		Role:         domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "passworduser",
		},
	}
	err := s.repo.Create(ctx, account)
	s.Require().NoError(err)
	s.Require().NotEmpty(account.ID)

	s.Run("Success", func() {
		newPassword := "new_hash"
		err := s.repo.UpdatePassword(ctx, account.ID, newPassword)
		s.NoError(err)

		// Verify in DB
		var model AccountModel
		objID, _ := primitive.ObjectIDFromHex(account.ID)
		err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
		s.NoError(err)
		s.Equal(newPassword, model.PasswordHash)
	})

	s.Run("Failure - Invalid ID", func() {
		err := s.repo.UpdatePassword(ctx, "badid", "irrelevant")
		s.ErrorIs(err, domain.ErrNotFound)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		err := s.repo.UpdatePassword(ctx, nonExistentID, "irrelevant")
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *AccountRepositoryTestSuite) TestExistsByEmail() {
	ctx := context.Background()

	// Arrange: Create a user with a known email
	account := &domain.Account{
		Name:         "Email Test User",
		Email:        "exists@example.com",
		PasswordHash: "hash",
		Role:         domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "emailtestuser",
		},
	}
	err := s.repo.Create(ctx, account)
	s.Require().NoError(err)
	s.Require().NotEmpty(account.ID)

	s.Run("Email Exists", func() {
		exists, err := s.repo.ExistsByEmail(ctx, "exists@example.com", "")
		s.NoError(err)
		s.True(exists)
	})

	s.Run("Email Does Not Exist", func() {
		exists, err := s.repo.ExistsByEmail(ctx, "nonexistent@example.com", "")
		s.NoError(err)
		s.False(exists)
	})

	s.Run("Email Exists But Excluded By ID", func() {
		exists, err := s.repo.ExistsByEmail(ctx, "exists@example.com", account.ID)
		s.NoError(err)
		s.False(exists)
	})

	s.Run("Email Exists For Different User", func() {
		// Create another user with different email
		otherAccount := &domain.Account{
			Name:         "Other User",
			Email:        "other@example.com",
			PasswordHash: "hash",
			Role:         domain.RoleUser,
			UserDetail: &domain.UserDetail{
				Username: "otheruser",
			},
		}
		err := s.repo.Create(ctx, otherAccount)
		s.Require().NoError(err)

		// Check if first user's email exists when excluding second user
		exists, err := s.repo.ExistsByEmail(ctx, "exists@example.com", otherAccount.ID)
		s.NoError(err)
		s.True(exists)
	})

	s.Run("Invalid Exclude ID", func() {
		_, err := s.repo.ExistsByEmail(ctx, "exists@example.com", "invalid-id")
		s.Error(err)
	})
}

func (s *AccountRepositoryTestSuite) TestExistsByUsername() {
	ctx := context.Background()

	// Arrange: Create a user with a known username
	account := &domain.Account{
		Name:         "Username Test User",
		Email:        "username@example.com",
		PasswordHash: "hash",
		Role:         domain.RoleUser,
		UserDetail: &domain.UserDetail{
			Username: "existinguser",
		},
	}
	err := s.repo.Create(ctx, account)
	s.Require().NoError(err)
	s.Require().NotEmpty(account.ID)

	s.Run("Username Exists", func() {
		exists, err := s.repo.ExistsByUsername(ctx, "existinguser", "")
		s.NoError(err)
		s.True(exists)
	})

	s.Run("Username Does Not Exist", func() {
		exists, err := s.repo.ExistsByUsername(ctx, "nonexistentuser", "")
		s.NoError(err)
		s.False(exists)
	})

	s.Run("Username Exists But Excluded By ID", func() {
		exists, err := s.repo.ExistsByUsername(ctx, "existinguser", account.ID)
		s.NoError(err)
		s.False(exists)
	})

	s.Run("Username Exists For Different User", func() {
		// Create another user with different username
		otherAccount := &domain.Account{
			Name:         "Other Username User",
			Email:        "otherusername@example.com",
			PasswordHash: "hash",
			Role:         domain.RoleUser,
			UserDetail: &domain.UserDetail{
				Username: "otherusername",
			},
		}
		err := s.repo.Create(ctx, otherAccount)
		s.Require().NoError(err)

		// Check if first user's username exists when excluding second user
		exists, err := s.repo.ExistsByUsername(ctx, "existinguser", otherAccount.ID)
		s.NoError(err)
		s.True(exists)
	})

	s.Run("Invalid Exclude ID", func() {
		_, err := s.repo.ExistsByUsername(ctx, "existinguser", "invalid-id")
		s.Error(err)
	})
}
