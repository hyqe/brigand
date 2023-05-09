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
		t.Skipf("missing env: MONGO")
	}

	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	defer mongoClient.Disconnect(ctx)

	coll := mongoClient.Database("brigand").Collection("metadata")
	doc := storage.Metadata{Id: "none", FileName: "none", CreateAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	mdc := storage.NewMongoMetadataClient(mongoClient)
	md, err := mdc.GetById(ctx, "none")
	require.NoError(t, err)

	require.Equal(t, md.Id, "none")
}
    
    
func Test_mongoMetadataClient_DeleteById_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("Missing env: MONGO")
	}

	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	defer mongoClient.Disconnect(ctx)

	coll := mongoClient.Database("brigand").Collection("metadata")
	doc := storage.Metadata{Id: "none", FileName: "none", CreateAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	mdc := storage.NewMongoMetadataClient(mongoClient)
	require.NoError(t, mdc.DeleteById(ctx, "none"))
}