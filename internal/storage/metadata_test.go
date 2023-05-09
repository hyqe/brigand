package storage_test

import (
	"context"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

func Test_mongoMetadataClient_DeleteById_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("Missing env: %s", MONGO)
	}

	// Open Client
	opts := options.Client().ApplyURI(MONGO)
	client, err := mongo.NewClient(opts)
	require.NoError(t, err)

	// Set Timed Connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	require.NoError(t, err)

	// Insert a Document to delete for the test
	coll := client.Database("brigand").Collection("metadata")
	doc := storage.Metadata{Id: "none", FileName: "none", CreateAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	// Delete the Document
	mdc := storage.NewMongoMetadataClient(client)
	require.NoError(t, mdc.DeleteById(ctx, "none"))
}
