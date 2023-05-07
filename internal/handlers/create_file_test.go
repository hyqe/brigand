package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewCreateFile_happy_path(t *testing.T) {
	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files", bytes.NewBuffer([]byte("string")))
	handlers.NewCreateFile(mockMetadataClient).ServeHTTP(w, r)

	// check status
	require.Equal(t, http.StatusOK, w.Code)

	// check id response
	mappy := make(map[string]any)
	json.NewDecoder(w.Body).Decode(&mappy)
	require.NotEmpty(t, mappy["id"])
}

func TestNewCreateFile_no_file(t *testing.T) {
	mockMetadataClient := &storage.MockMetadataClient{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files", nil)
	handlers.NewCreateFile(mockMetadataClient).ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
