package repository_test

// import (
// 	"context"
// 	"testing"

// 	"EthioGuide/domain"
// 	"EthioGuide/repository"

// 	"github.com/stretchr/testify/suite"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
// )

// type PreferencesRepositorySuite struct {
//     suite.Suite
//     mt   *mtest.T
//     repo domain.IPreferencesRepository
// }

// func (s *PreferencesRepositorySuite) SetupTest() {
//     s.mt = mtest.New(s.T(), mtest.NewOptions().ClientType(mtest.Mock))
//     s.repo = repository.NewPreferencesRepositoryMongo(s.mt.DB)
// }

// func (s *PreferencesRepositorySuite) TearDownTest() {
//     s.mt.Close()
// }

// func (s *PreferencesRepositorySuite) TestCreateAndGetByUserID() {
//     s.mt.Run("Create and GetByUserID", func(mt *mtest.T) {
//         ctx := context.Background()
//         userID := "user123"
//         pref := &domain.Preferences{
//             UserID:            userID,
//             PreferredLang:     "en",
//             PushNotification:  true,
//             EmailNotification: false,
//         }

//         mt.AddMockResponses(mtest.CreateSuccessResponse())
//         err := s.repo.Create(ctx, pref)
//         s.NoError(err)
//         s.NotEmpty(pref.ID)

//         // Prepare a mock response for FindOne
//         oid, _ := primitive.ObjectIDFromHex(pref.ID)
//         mt.AddMockResponses(mtest.CreateCursorResponse(
//             1, "preferences.test", mtest.FirstBatch,
//             bson.D{
//                 {"_id", oid},
//                 {"user_id", userID},
//                 {"preferredLang", "en"},
//                 {"pushNotification", true},
//                 {"emailNotification", false},
//             },
//         ))

//         got, err := s.repo.GetByUserID(ctx, userID)
//         s.NoError(err)
//         s.Equal(userID, got.UserID)
//         s.Equal("en", string(got.PreferredLang))
//         s.Equal(true, got.PushNotification)
//         s.Equal(false, got.EmailNotification)
//     })
// }

// func (s *PreferencesRepositorySuite) TestUpdateByUserID() {
//     s.mt.Run("UpdateByUserID", func(mt *mtest.T) {
//         ctx := context.Background()
//         userID := "user456"
//         pref := &domain.Preferences{
//             UserID:            userID,
//             PreferredLang:     "am",
//             PushNotification:  false,
//             EmailNotification: true,
//         }

//         mt.AddMockResponses(mtest.CreateSuccessResponse())
//         err := s.repo.UpdateByUserID(ctx, userID, pref)
//         s.NoError(err)
//     })
// }

// func TestPreferencesRepositorySuite(t *testing.T) {
//     suite.Run(t, new(PreferencesRepositorySuite))
// }