package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChecklistItemModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserProcedureID primitive.ObjectID `bson:"user_procedure_id"`
	Type            string             `bson:"type"`
	Content         string             `bson:"content"`
	IsChecked       bool               `bson:"is_checked"`
}

type ChecklistRepository struct {
	collectionProcedure     *mongo.Collection
	collectionChecklist     *mongo.Collection
	collectionUserProcedure *mongo.Collection
}

func TodomainChecklist(check *ChecklistItemModel) *domain.Checklist {
	return &domain.Checklist{
		ID:              check.ID.Hex(),
		UserProcedureID: check.UserProcedureID.Hex(),
		Type:            check.Type,
		Content:         check.Content,
		IsChecked:       check.IsChecked,
	}
}

func NewChecklistRepository(db *mongo.Database) *ChecklistRepository {
	collcheck := db.Collection("checklist")
	collprocedure := db.Collection("procedures")
	colluserprocdr := db.Collection("user_procedure")
	return &ChecklistRepository{
		collectionChecklist:     collcheck,
		collectionUserProcedure: colluserprocdr,
		collectionProcedure:     collprocedure,
	}
}

func (cr *ChecklistRepository) CreateChecklist(ctx context.Context, userid, procdureID string) (*domain.UserProcedure, error) {
	objID, err := primitive.ObjectIDFromHex(procdureID)
	if err != nil {
		return nil, err
	}

	objuserID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return nil, err
	}

	var procedure ProcedureModel
	if err := cr.collectionProcedure.FindOne(ctx, bson.M{"_id": objID}).Decode(&procedure); err != nil {
		return nil, err
	}

	doc := &UserProcedureModel{
		UserID:      objuserID,
		ProcedureID: objID,
		Percent:     0,
		UpdatedAt:   time.Now(),
	}

	res, err := cr.collectionUserProcedure.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	userprocedureID := res.InsertedID.(primitive.ObjectID)
	documents := make([]ChecklistItemModel, 0, (len(procedure.Content.Prerequisites) + len(procedure.Content.Steps)))
	for _, req := range procedure.Content.Prerequisites {
		documents = append(documents, ChecklistItemModel{
			UserProcedureID: userprocedureID,
			Type:            "Requirement",
			Content:         req,
			IsChecked:       false,
		})
	}

	for _, step := range procedure.Content.Steps {
		documents = append(documents, ChecklistItemModel{
			UserProcedureID: userprocedureID,
			Type:            "Step",
			Content:         step,
			IsChecked:       false,
		})
	}

	if _, err := cr.collectionChecklist.InsertMany(ctx, toInsertSlice(documents)); err != nil {
		return nil, err
	}

	doc.ID = userprocedureID
	return ToDomainUserProcedure(doc), nil
}

func (cr *ChecklistRepository) GetProcedures(ctx context.Context, userid string) ([]*domain.UserProcedure, error) {
	objID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": objID}
	cursor, err := cr.collectionUserProcedure.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var procedures []*domain.UserProcedure
	for cursor.Next(ctx) {
		var procedure UserProcedureModel
		if err := cursor.Decode(&procedure); err != nil {
			return nil, err
		}
		procedures = append(procedures, ToDomainUserProcedure(&procedure))
	}

	return procedures, nil
}

func (cr *ChecklistRepository) GetChecklistByUserProcedureID(ctx context.Context, userprocedureID string) ([]*domain.Checklist, error) {
	objID, err := primitive.ObjectIDFromHex(userprocedureID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_procedure_id": objID}
	cursor, err := cr.collectionChecklist.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var checklists []*domain.Checklist
	for cursor.Next(ctx) {
		var checklistmodel ChecklistItemModel
		if err := cursor.Decode(&checklistmodel); err != nil {
			return nil, err
		}

		checklists = append(checklists, TodomainChecklist(&checklistmodel))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return checklists, nil
}

func (cr *ChecklistRepository) ToggleCheck(ctx context.Context, checklistID string) error {
	objID, err := primitive.ObjectIDFromHex(checklistID)
	if err != nil {
		return err // Invalid ID format
	}

	// Use a session for the transaction
	session, err := cr.collectionChecklist.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// The WithTransaction function handles starting, committing, and aborting the transaction.
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// --- Step 1: Atomically toggle the item and get its updated state ---
		// We use FindOneAndUpdate to get the document back after the update.
		// We need the `user_procedure_id` from it.
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		filter := bson.M{"_id": objID}
		// The update pipeline '$not' efficiently flips the boolean value.
		update := bson.A{bson.D{{Key: "$set", Value: bson.D{{Key: "is_checked", Value: bson.D{{Key: "$not", Value: "$is_checked"}}}}}}}

		var updatedChecklistItem ChecklistItemModel
		err := cr.collectionChecklist.FindOneAndUpdate(sessCtx, filter, update, opts).Decode(&updatedChecklistItem)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, domain.ErrNotFound
			}
			return nil, err // Other database error
		}

		userProcedureID := updatedChecklistItem.UserProcedureID

		// --- Step 2: Count total and checked items for the parent procedure ---
		// Both counts must be performed within the session context.
		totalItemsFilter := bson.M{"user_procedure_id": userProcedureID}
		totalCount, err := cr.collectionChecklist.CountDocuments(sessCtx, totalItemsFilter)
		if err != nil {
			return nil, err
		}

		if totalCount == 0 {
			// This is a safety check; should not happen if a checklist item exists.
			return nil, domain.ErrNotFound
		}

		checkedItemsFilter := bson.M{"user_procedure_id": userProcedureID, "is_checked": true}
		checkedCount, err := cr.collectionChecklist.CountDocuments(sessCtx, checkedItemsFilter)
		if err != nil {
			return nil, err
		}

		// --- Step 3: Calculate the new percentage ---
		// Use float64 for division to get an accurate ratio.
		percentage := int((float64(checkedCount) / float64(totalCount)) * 100.0)
		var status string
		if percentage == 100 {
			status = "Completed"
		} else if percentage > 0 {
			status = "In Progress"
		} else {
			status = "Not Started"
		}

		// --- Step 4: Update the parent UserProcedure document ---
		parentUpdateFilter := bson.M{"_id": userProcedureID}
		parentUpdate := bson.M{
			"$set": bson.M{
				"percent":    percentage,
				"status":     status,
				"updated_at": time.Now(),
			},
		}

		result, err := cr.collectionUserProcedure.UpdateOne(sessCtx, parentUpdateFilter, parentUpdate)
		if err != nil {
			return nil, err
		}
		if result.MatchedCount == 0 {
			return nil, domain.ErrNotFound
		}

		// If everything succeeded, return nil to commit the transaction.
		return nil, nil
	})

	return err
}

func (cr *ChecklistRepository) FindCheck(ctx context.Context, checklistID string) (*domain.Checklist, error) {
	objID, err := primitive.ObjectIDFromHex(checklistID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	var checklist ChecklistItemModel
	if err := cr.collectionChecklist.FindOne(ctx, filter).Decode(&checklist); err != nil {
		return nil, err
	}

	return TodomainChecklist(&checklist), nil
}

func (cr *ChecklistRepository) CountDocumentsChecklist(ctx context.Context, filter interface{}) (int64, error) {
	countDoc, err := cr.collectionChecklist.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return countDoc, nil
}

func (cr *ChecklistRepository) UpdateUserProcedure(ctx context.Context, filter interface{}, update map[string]interface{}) error {
	updateDoc := bson.M{
		"$set": update,
	}
	if _, err := cr.collectionUserProcedure.UpdateOne(ctx, filter, updateDoc); err != nil {
		return err
	}

	return nil
}

func toInsertSlice(items []ChecklistItemModel) []interface{} {
	result := make([]interface{}, len(items))
	for i, v := range items {
		result[i] = v
	}

	return result
}
