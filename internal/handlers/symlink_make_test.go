package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	// "fmt"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func symlinkParams(r *http.Request) map[string]string {

	mappy := make(map[string]string)
	mappy["hash"] = r.URL.Query().Get("hash")
	mappy["expiration"] = r.URL.Query().Get("expiration")
	mappy["id"] = r.URL.Query().Get("id")
	mappy["name"] = r.URL.Query().Get("name")

	return mappy
}

func Test_MakeSymlink_happy_path(t *testing.T) {
	fc := []byte("FileContent")
	fileContent := bytes.NewBuffer(fc)

	hmacSecret := "MySecret"
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
		return r.URL.Query().Get("name")
	}

	path := "/symlink/make?name=Breaker"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, path, fileContent)

	handlers.MakeSymlink(mockMetadataClient, mockFileUploader, getFileId, hmacSecret).ServeHTTP(w, r)

	var symlink handlers.Symlink
	err := json.NewDecoder(w.Body).Decode(&symlink)
	require.NoError(t, err)

	// Check if it gets the correct amount of params
	params := symlinkParams(httptest.NewRequest(http.MethodGet, symlink.Link, nil))
	require.GreaterOrEqual(t, len(params), 4)

	// Check if the link originated from us by comparing the hashes
	require.True(t, handlers.CheckHash(params, hmacSecret))

	require.Equal(t, http.StatusOK, w.Code)
}

func Test_MakeSymlink_no_file_content(t *testing.T) {
	hmacSecret := "MySecret"
	mockMetadataClient := &storage.MockMetadataClient{
		CreateFunc: func(ctx context.Context, md *storage.Metadata) error {
			return nil
		},
	}

	// mockFileUploader := func(file io.Reader, filename string) error {
	// // require.Equal(t, nil, file)
	// // return fmt.Errorf("")
	// return nil
	// }

	getFileId := func(r *http.Request) string {
		return r.URL.Query().Get("name")
	}

	path := "/symlink/make?name=Breaker"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, path, nil)

	handlers.MakeSymlink(mockMetadataClient, nil, getFileId, hmacSecret).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)

}
