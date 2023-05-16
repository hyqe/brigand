package handlers_test

import (
	"context"
	"io"

	"github.com/hyqe/brigand/internal/server"
	"github.com/hyqe/brigand/internal/storage"

	"fmt"
	"net/http"
	"net/http/httptest"

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

	//mockGetRequestFileId := func(r *http.Request) string {
	//	return md.Id
	//}
	//handler := handlers.NewGetFileById(mockMetadataClient, mockFileDownloader, mockGetRequestFileId)

	srv := httptest.NewServer(server.Routes(mockMetadataClient, mockFileDownloader))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/files/%v", srv.URL, md.Id), nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(b))
}

//func Test_GetFileById_bad_file_id(t *testing.T) {
//	mock := &storage.MockMetadataClient{}
//	mockFileDownloader := storage.MockNewS3FileDownloader(fmt.Errorf("NoFileFound"))
//
//	w := httptest.NewRecorder()
//	r := httptest.NewRequest(http.MethodGet, "/files/{fileId}", nil)
//	r = mux.SetURLVars(r, map[string]string{"fileId": uuid.New().String()})
//	handlers.NewGetFileById(mock, mockFileDownloader).ServeHTTP(w, r)
//
//	require.Equal(t, 404, w.Result().StatusCode)
//}
