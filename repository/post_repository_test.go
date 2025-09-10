package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.IPostRepository
	collection *mongo.Collection
}

func (s *PostRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewPostRepository(s.db)
	s.collection = s.db.Collection("posts")
}

func (s *PostRepositoryTestSuite) TearDownSuite() {
	s.db.Drop(context.Background())
}

func (s *PostRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	s.collection.DeleteMany(context.Background(), bson.M{})
}

func TestPostRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (s *PostRepositoryTestSuite) TestCreatePost() {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	procID := primitive.NewObjectID()

	post := &domain.Post{
		UserID:     userID.Hex(),
		Title:      "New Post Title",
		Content:    "Post content here.",
		Procedures: []string{procID.Hex()},
		Tags:       []string{"testing"},
	}

	createdPost, err := s.repo.CreatePost(ctx, post)
	s.NoError(err)
	s.NotNil(createdPost)
	s.NotEmpty(createdPost.ID)

	// Verify in DB
	var model PostModel
	objID, _ := primitive.ObjectIDFromHex(createdPost.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.NoError(err)
	s.Equal("New Post Title", model.Title)
	s.Equal(userID, model.UserID)
	s.Len(model.ProcedureIDs, 1)
	s.Equal(procID, model.ProcedureIDs[0])
}

func (s *PostRepositoryTestSuite) TestGetPostByID() {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	// Arrange
	model := &PostModel{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Title:  "Find Me",
	}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	s.Run("Success", func() {
		post, err := s.repo.GetPostByID(ctx, model.ID.Hex())
		s.NoError(err)
		s.NotNil(post)
		s.Equal("Find Me", post.Title)
	})

	s.Run("Failure - Not Found", func() {
		_, err := s.repo.GetPostByID(ctx, primitive.NewObjectID().Hex())
		s.ErrorIs(err, domain.ErrPostNotFound)
	})
}

func (s *PostRepositoryTestSuite) TestGetPosts() {
	ctx := context.Background()
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()
	procID1 := primitive.NewObjectID()
	procID2 := primitive.NewObjectID()

	// Arrange
	posts := []interface{}{
		PostModel{UserID: userID1, Title: "Alpha Post", Tags: []string{"go", "testing"}, ProcedureIDs: []primitive.ObjectID{procID1}},
		PostModel{UserID: userID2, Title: "Beta Post", Tags: []string{"go", "tutorial"}, ProcedureIDs: []primitive.ObjectID{procID2}},
		PostModel{UserID: userID1, Title: "Gamma Test", Tags: []string{"testing"}, ProcedureIDs: []primitive.ObjectID{procID1, procID2}},
	}
	_, err := s.collection.InsertMany(ctx, posts)
	s.Require().NoError(err)

	s.Run("Filter by Title", func() {
		title := "Alpha"
		opts := domain.PostFilters{Title: &title}
		results, total, err := s.repo.GetPosts(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Len(results, 1)
	})

	s.Run("Filter by UserID", func() {
		userIDHex := userID1.Hex()
		opts := domain.PostFilters{UserId: &userIDHex}
		results, total, err := s.repo.GetPosts(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total)
		s.Len(results, 2)
	})

	s.Run("Filter by Tags ($all)", func() {
		opts := domain.PostFilters{Tags: []string{"go", "testing"}}
		results, total, err := s.repo.GetPosts(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Len(results, 1)
		s.Equal("Alpha Post", results[0].Title)
	})

	s.Run("Filter by ProcedureIDs ($all)", func() {
		opts := domain.PostFilters{ProcedureID: []string{procID1.Hex(), procID2.Hex()}}
		results, total, err := s.repo.GetPosts(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Len(results, 1)
		s.Equal("Gamma Test", results[0].Title)
	})
}

func (s *PostRepositoryTestSuite) TestUpdatePost() {
	ctx := context.Background()
	authorID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	// Arrange
	model := &PostModel{
		ID:     primitive.NewObjectID(),
		UserID: authorID,
		Title:  "Original Title",
	}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	s.Run("Success - Author updates own post", func() {
		postUpdate := &domain.Post{
			ID:      model.ID.Hex(),
			UserID:  authorID.Hex(), // Correct author
			Title:   "Updated Title",
			Content: "Updated Content",
		}
		updatedPost, err := s.repo.UpdatePost(ctx, postUpdate)
		s.NoError(err)
		s.NotNil(updatedPost)
		s.Equal("Updated Title", updatedPost.Title)
	})

	s.Run("Failure - Another user tries to update", func() {
		postUpdate := &domain.Post{
			ID:     model.ID.Hex(),
			UserID: otherUserID.Hex(), // Incorrect author
			Title:  "Should Fail",
		}
		_, err := s.repo.UpdatePost(ctx, postUpdate)
		s.ErrorIs(err, domain.ErrPermissionDenied)
	})
}

func (s *PostRepositoryTestSuite) TestDeletePost() {
	ctx := context.Background()
	authorID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	adminID := primitive.NewObjectID()

	s.Run("Success - Author deletes own post", func() {
		// Arrange
		s.BeforeTest("", "") // Clean collection
		model := &PostModel{ID: primitive.NewObjectID(), UserID: authorID}
		_, err := s.collection.InsertOne(ctx, model)
		s.Require().NoError(err)

		// Act
		err = s.repo.DeletePost(ctx, model.ID.Hex(), authorID.Hex(), string(domain.RoleUser))
		s.NoError(err)

		// Assert
		count, _ := s.collection.CountDocuments(ctx, bson.M{"_id": model.ID})
		s.Equal(int64(0), count)
	})

	s.Run("Success - Admin deletes post", func() {
		// Arrange
		s.BeforeTest("", "")
		model := &PostModel{ID: primitive.NewObjectID(), UserID: authorID}
		_, err := s.collection.InsertOne(ctx, model)
		s.Require().NoError(err)

		// Act
		err = s.repo.DeletePost(ctx, model.ID.Hex(), adminID.Hex(), string(domain.RoleAdmin))
		s.NoError(err)

		// Assert
		count, _ := s.collection.CountDocuments(ctx, bson.M{"_id": model.ID})
		s.Equal(int64(0), count)
	})

	s.Run("Failure - Another user tries to delete", func() {
		// Arrange
		s.BeforeTest("", "")
		model := &PostModel{ID: primitive.NewObjectID(), UserID: authorID}
		_, err := s.collection.InsertOne(ctx, model)
		s.Require().NoError(err)

		// Act
		err = s.repo.DeletePost(ctx, model.ID.Hex(), otherUserID.Hex(), string(domain.RoleUser))
		s.ErrorIs(err, domain.ErrPermissionDenied)

		// Assert
		count, _ := s.collection.CountDocuments(ctx, bson.M{"_id": model.ID})
		s.Equal(int64(1), count) // Ensure it was not deleted
	})
}
