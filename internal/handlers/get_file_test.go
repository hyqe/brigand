package handlers_test

import (
	"github.com/google/uuid"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"

	"testing"
)

func Test_GetFile_happy_path(t *testing.T) {
	md := storage.NewMetadata("Breaker")
	mockMetadataClient := &storage.MockMetadataClient{}
	mockFileDownloader := storage.MockNewS3FileDownloader(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
	r = mux.SetURLVars(r, map[string]string{"fileId": md.Id})
	handlers.NewGetFileById(mockMetadataClient, mockFileDownloader).ServeHTTP(w, r)

	if w.Result().StatusCode < 200 || w.Result().StatusCode > 299 {
		require.True(t, false)
	}
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
