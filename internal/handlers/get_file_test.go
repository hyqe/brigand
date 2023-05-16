package handlers_test

import (
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"

	"context"
	"testing"
)

func Test_GetFile_happy_path(t *testing.T) {
	md := storage.NewMetadata("Breaker")
	fileContent := "FileContent"
	mockMetadataClient := &storage.MockMetadataClient{
		GetFileById: func(ctx context.Context, id string) (*storage.Metadata, error) {
			require.Equal(t, md.Id, id)
			return md, nil
		},
	}
	mockFileDownloader := func(file io.Writer, filename string) error {
		require.Equal(t, md.Id, filename)
		fmt.Fprint(file, fileContent)
		return nil
	}
	mockGetFileId := func(r *http.Request) string {
		return mux.Vars(r)["fileId"]
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	r = mux.SetURLVars(r, map[string]string{"fileId": md.Id})
	handlers.NewGetFileById(mockMetadataClient, mockFileDownloader, mockGetFileId).ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, fileContent, w.Body.String())
}

func Test_GetFileById_bad_file_id(t *testing.T) {
	md := storage.NewMetadata("Breaker")
	mockMetaDataClient := &storage.MockMetadataClient{
		GetFileById: func(ctx context.Context, id string) (*storage.Metadata, error) {
			return md, fmt.Errorf("NoFileId")
		},
	}
	mockFileDownloader := func(file io.Writer, filename string) error {
		return nil
	}
	mockGetFileId := func(r *http.Request) string {
		return ""
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	r = mux.SetURLVars(r, map[string]string{"fileId": "badid"})
	handlers.NewGetFileById(mockMetaDataClient, mockFileDownloader, mockGetFileId).ServeHTTP(w, r)

	require.Equal(t, 404, w.Result().StatusCode)
}
