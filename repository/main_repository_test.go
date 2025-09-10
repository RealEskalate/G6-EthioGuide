package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// testDBClient will be a global variable accessible to all tests in the package.
var testDBClient *mongo.Client

const testDBName = "ethio-guide-test"

// TestMain is the special setup/teardown function for the entire package.
func TestMain(m *testing.M) {
	// 1. Setup Phase
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start a MongoDB container
	mongodbContainer, err := mongodb.Run(ctx,
		"mongo:6.0",
		mongodb.WithReplicaSet("rs0"),
	)
	if err != nil {
		log.Fatalf("Failed to start MongoDB container: %s", err)
	}

	// Get the connection string for the container
	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("Failed to get MongoDB connection string: %s", err)
	}

	// Create a new MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to test MongoDB: %s", err)
	}

	// Ping the database to ensure connection is live
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping test MongoDB: %s", err)
	}

	// Assign the client to our global variable
	testDBClient = client
	fmt.Println("✅ Test MongoDB container is up and running.")

	// 2. Run the tests
	// os.Exit() does not respect defers, so we wrap the teardown in a function.
	exitCode := m.Run()

	// 3. Teardown Phase
	fmt.Println("🔴 Tearing down test MongoDB container...")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("Failed to disconnect from test MongoDB: %s", err)
	}
	if err := mongodbContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("Failed to terminate MongoDB container: %s", err)
	}

	// Exit with the result of the tests
	os.Exit(exitCode)
}
