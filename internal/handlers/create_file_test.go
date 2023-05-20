package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewCreateFile_happy_path(t *testing.T) {

	var MD *storage.Metadata

	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			MD = md
			return nil
		},
	}

	mockFileUploader := func(file io.Reader, filename string) error {
		require.Equal(t, MD.Id, filename)
		return nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files?name=foo", bytes.NewBuffer([]byte("my file content")))
	handlers.NewCreateFile(mockMetadataClient, mockFileUploader).ServeHTTP(w, r)

	// check status
	require.Equal(t, http.StatusOK, w.Code)

	// check id response
	var resp storage.Metadata
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, MD.Id, resp.Id)
}

func TestNewCreateFile_reject_empty_body(t *testing.T) {
	mockFileUploader := func(file io.Reader, filename string) error {
		return fmt.Errorf("no body content")
	}

	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files", nil)
	handlers.NewCreateFile(mockMetadataClient, mockFileUploader).ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
