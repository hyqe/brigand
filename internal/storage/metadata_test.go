package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/hyqe/brigand/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/require"
)

func Test_mongoMetadataClient_GetById_happy_path(t *testing.T) {
	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("missing env: %s", "MONGO")
	}

	// Get the client
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO))
	require.NoError(t, err)

	ctx := context.Background()
	defer client.Disconnect(ctx)

	mdc := storage.NewMongoMetadataClient(client)

	// There could or could not be a value?
	_, err = mdc.GetById(ctx, "magicalid")
	require.NoError(t, err)

}
