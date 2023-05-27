package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MakeSymlink_happy_path(t *testing.T) {
	fc := []byte("FileContent")
	fileContent := bytes.NewBuffer(fc)

	symlinkSecret := "MySecret"
	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			return nil
		},
	}

	mockFileUploader := func(file io.Reader, filename string) error {
		byties, err := io.ReadAll(file)
		require.NoError(t, err)

		for key, bytie := range byties {
			require.Equal(t, fc[key], bytie)
		}
		require.Equal(t, fc, byties)

		return nil
	}

	getFileId := func(r *http.Request) string {
		return "breaker"
	}

	path := "/"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, path, fileContent)

	handlers.MakeSymlink(mockMetadataClient, mockFileUploader, getFileId, symlinkSecret).ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	handlerGeneratedSymlink := &storage.Symlink{}
	err := json.NewDecoder(w.Body).Decode(handlerGeneratedSymlink)
	require.NoError(t, err)

	fakeR := httptest.NewRequest(http.MethodPost, handlerGeneratedSymlink.Link, nil)
	symlink := storage.SymlinkFromQuery(fakeR)
	symlink.Name = "breaker"

	require.True(t, handlers.CheckSymlinkHash(symlink, symlinkSecret))
	require.Equal(t, http.StatusOK, w.Code)
}

func Test_MakeSymlink_no_file_content(t *testing.T) {
	symlinkSecret := "MySecret"
	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			return nil
		},
	}

	mockFileUploader := func(file io.Reader, filename string) error {
		return nil
	}

	getFileId := func(r *http.Request) string {
		return "breaker"
	}

	path := "/symlink/make?name=Breaker"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, path, nil)

	handlers.MakeSymlink(mockMetadataClient, mockFileUploader, getFileId, symlinkSecret).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)

}
