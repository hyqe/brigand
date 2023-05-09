package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/hyqe/brigand/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stretchr/testify/require"
	"time"
)

func Test_mongoMetadataClient_GetById_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("missing env: %s", "MONGO")
	}

	// Open Client
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO))
	require.NoError(t, err)

	// Set Timed Connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	require.NoError(t, err)

	// Insert Document for test
	coll := client.Database("brigand").Collection("metadata")
	doc := storage.Metadata{Id: "none", FileName: "none", CreateAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	// Get Document
	mdc := storage.NewMongoMetadataClient(client)
	_, err = mdc.GetById(ctx, "none")
	require.NoError(t, err)

}
