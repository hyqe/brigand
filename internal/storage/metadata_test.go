package storage_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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
	doc := storage.Metadata{Id: "none", FileName: "none", CreatedAt: time.Now()}
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
	doc := storage.Metadata{Id: "none", FileName: "none", CreatedAt: time.Now()}
	_, err = coll.InsertOne(context.TODO(), doc)
	require.NoError(t, err)

	mdc := storage.NewMongoMetadataClient(mongoClient)
	require.NoError(t, mdc.DeleteById(ctx, "none"))
}

func Test_CreateMetadataIndex_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("Missing env: MONGO")
	}

	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	require.NoError(t, err)
	defer mongoClient.Disconnect(ctx)

	collName := uuid.New().String()
	coll := mongoClient.Database("brigand").Collection(collName)
	defer coll.Drop(ctx)

	err = storage.CreateMetadataIndex(ctx, coll)
	require.NoError(t, err)

	md := storage.NewMetadata("HappyPath")
	_, err = coll.InsertOne(context.TODO(), md)
	require.NoError(t, err)

	md = &storage.Metadata{Id: md.Id, FileName: "HappyPath2", CreatedAt: time.Now().UTC()}
	_, err = coll.InsertOne(context.TODO(), md)

	require.Error(t, err)

}
