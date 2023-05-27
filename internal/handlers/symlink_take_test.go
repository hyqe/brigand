package handlers_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
)

func formatSymlink(md *storage.Metadata, hmacSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(time.Hour).Format(time.RFC3339)
	signature := storage.MakeSymlinkSignature(theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(hmacSecret))
	_, err := h.Write([]byte(signature))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	query := storage.MakeSymlinkQuery(theTime, md.Id)
	symlink := fmt.Sprintf("/symlink/take/{name}?hash=%s&%s", hash, query)

	return symlink
}

func formatSymlinkWithBadTime(md *storage.Metadata, hmacSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(-time.Hour).Format(time.RFC3339)
	signature := storage.MakeSymlinkSignature(theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(hmacSecret))
	_, err := h.Write([]byte(signature))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	query := storage.MakeSymlinkQuery(theTime, md.Id)
	symlink := fmt.Sprintf("/symlink/take/{name}?hash=%s&%s", hash, query)

	return symlink
}

func Test_TakeSymlink_happy_path(t *testing.T) {
	fileContent := bytes.NewBuffer([]byte("FileContent"))

	md := storage.NewMetadata("Breaker")
	symlinkSecret := "MySecret"
	mockSymlink := formatSymlink(md, symlinkSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		_, err := io.Copy(file, fileContent)
		require.NoError(t, err)

		return nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)

	r = mux.SetURLVars(r, map[string]string{"name": md.FileName})

	handlers.TakeSymlink(fileDownloader, symlinkSecret, storage.SymlinkFromQuery).ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	// LOL why cant i copy this?
	fileContent = bytes.NewBuffer([]byte("FileContent"))
	require.Equal(t, fileContent, w.Body)

}

func Test_TakeSymlink_unauthentic_symlink(t *testing.T) {

	md := storage.NewMetadata("Breaker")
	dirtySymlinkSecret := "DirtySecret"
	mockSymlink := formatSymlink(md, dirtySymlinkSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		return nil
	}

	realSymlinkSecret := "RealSecret"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)

	r = mux.SetURLVars(r, map[string]string{"name": md.FileName})

	handlers.TakeSymlink(fileDownloader, realSymlinkSecret, storage.SymlinkFromQuery).ServeHTTP(w, r)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func Test_TakeSymlink_expired_symlink(t *testing.T) {

	md := storage.NewMetadata("Breaker")
	symlinkSecret := "MySecret"
	mockSymlink := formatSymlinkWithBadTime(md, symlinkSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		return nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)

	r = mux.SetURLVars(r, map[string]string{"name": md.FileName})

	handlers.TakeSymlink(fileDownloader, symlinkSecret, storage.SymlinkFromQuery).ServeHTTP(w, r)

	require.Equal(t, http.StatusNotAcceptable, w.Code)
}
