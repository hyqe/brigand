package storage_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
)

func Test_mongoMetadataClient_GetById_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("missing env: %s", "MONGO")
	}

	// Open CLient
	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	mongoClient.Connect(ctx)
	defer mongoClient.Disconnect(ctx)

	// Insert Document for test
	coll := mongoClient.Database("brigand").Collection("metadata")
	doc := storage.Metadata{Id: "none", FileName: "none", CreateAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	// Get Document
	mdc := storage.NewMongoMetadataClient(mongoClient)
	md, err := mdc.GetById(ctx, "none")
	require.NoError(t, err)

	require.Equal(t, md.Id, "none")
}
