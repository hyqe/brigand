package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"encoding/json"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewGetFileById_happy_path(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/1234", nil)

	MONGO, ok := os.LookupEnv("MONGO")
	if !ok {
		t.Skipf("Missing env: MONGO")
	}

	// Open a Mongo Client
	ctx := context.Background()
	mongoClient, err := storage.NewMongoClient(ctx, MONGO)
	require.NoError(t, err)

	// Get the File
	mmdClient := storage.NewMongoMetadataClient(mongoClient)
	NewGetFileById(mmdClient).ServeHTTP(w, r)

	// Pretend this is a file deserialize/check
	var magic_data storage.Magicdata
	err = json.NewDecoder(w.Body).Decode(&magic_data)
	require.NoError(t, err)
}