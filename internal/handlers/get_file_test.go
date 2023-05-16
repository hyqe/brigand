package handlers_test

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"

	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"testing"
)

func Test_GetFile_happy_path(t *testing.T) {
	md := storage.NewMetadata("Breaker")
	fileContent := "this is some content"

	mockMetadataClient := &storage.MockMetadataClient{
		GetByIdFunc: func(ctx context.Context, id string) (*storage.Metadata, error) {
			require.Equal(t, md.Id, id)
			return md, nil
		},
	}

	mockFileDownloader := func(file io.Writer, filename string) error {
		require.Equal(t, md.Id, filename)
		fmt.Fprint(file, fileContent)
		return nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ANYTHING GOES", nil)

	mockGetFileId := func(r *http.Request) string { return md.Id }
	handlers.NewGetFileById(mockMetadataClient, mockFileDownloader, mockGetFileId).ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, fileContent, w.Body.String())
}

func Test_GetFileById_bad_file_id(t *testing.T) {
	mock := &storage.MockMetadataClient{}
	mockFileDownloader := storage.MockNewS3FileDownloader(fmt.Errorf("NoFileFound"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	r = mux.SetURLVars(r, map[string]string{"fileId": uuid.New().String()})
	handlers.NewGetFileById(mock, mockFileDownloader).ServeHTTP(w, r)

	require.Equal(t, 404, w.Result().StatusCode)
}
