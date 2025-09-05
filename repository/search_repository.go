package repository

import (
	"EthioGuide/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

type SearchRepository struct {
	db *mongo.Database
}

func NewSearchRepository(db *mongo.Database) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

// Main Search method
func (r *SearchRepository) Search(ctx context.Context, filters domain.SearchFilterRequest) (*domain.SearchResult, error) {
	// Use an errgroup to run the two queries in parallel for better performance.
	g, gCtx := errgroup.WithContext(ctx)

	var procedures []*domain.Procedure
	var organizations []*domain.AccountOrgSearch

	// Run findProcedures in a separate goroutine
	g.Go(func() error {
		var err error
		procedures, err = r.FindProcedures(gCtx, filters)
		return err
	})

	// Run findOrganizations in a separate goroutine
	g.Go(func() error {
		var err error
		organizations, err = r.FindOrganizations(gCtx, filters)
		return err
	})

	// Wait for both queries to complete.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Combine the results into the final struct.
	return &domain.SearchResult{
		Procedures:    procedures,
		Organizations: organizations,
	}, nil
}

// --- Helper Methods ---

// findProcedures performs a text search on the procedures collection.
func (r *SearchRepository) FindProcedures(ctx context.Context, filters domain.SearchFilterRequest) ([]*domain.Procedure, error) {
	// CRITICAL FIX: Removed the incorrect `{Key: "role", ...}` filter.
	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: filters.Query}}}}

	findOptions := options.Find()
	findOptions.SetLimit(filters.Limit)
	findOptions.SetSkip((filters.Page - 1) * filters.Limit)

	cursor, err := r.db.Collection("procedures").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Procedure
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// findOrganizations performs a text search on the accounts collection for organizations.
func (r *SearchRepository) FindOrganizations(ctx context.Context, filters domain.SearchFilterRequest) ([]*domain.AccountOrgSearch, error) {
	filter := bson.D{
		{Key: "role", Value: domain.RoleOrg},
		{Key: "$text", Value: bson.D{{Key: "$search", Value: filters.Query}}},
	}

	findOptions := options.Find()
	findOptions.SetLimit(filters.Limit)
	findOptions.SetSkip((filters.Page - 1) * filters.Limit)

	cursor, err := r.db.Collection("accounts").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accountModels []*AccountModel
	if err = cursor.All(ctx, &accountModels); err != nil {
		return nil, err
	}

	// CRITICAL FIX: Convert directly to the public-facing AccountOrgSearch struct.
	results := make([]*domain.AccountOrgSearch, 0, len(accountModels))
	for _, model := range accountModels {
		domainAccount := toDomainAccount(model)                   // First convert to domain model
		results = append(results, domain.ToSearch(domainAccount)) // Then to search result model
	}

	return results, nil
}
