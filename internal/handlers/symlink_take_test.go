package handlers_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func symlinkParamsx(r *http.Request) map[string]string {

	mappy := make(map[string]string)
	mappy["hash"] = r.URL.Query().Get("hash")
	mappy["expiration"] = r.URL.Query().Get("expiration")
	mappy["id"] = r.URL.Query().Get("id")
	mappy["name"] = r.URL.Query().Get("name")

	return mappy
}

func formatSymlink(md *storage.Metadata, hmacSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(time.Hour).Format(time.RFC3339)
	path := fmt.Sprintf("expiration=%s&id=%s&name=%s", theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(hmacSecret))
	_, err := h.Write([]byte(path))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	symlink := fmt.Sprintf("/symlink/take?hash=%s&%s", hash, path)

	return symlink
}

func formatSymlinkWithBadTime(md *storage.Metadata, hmacSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(-time.Hour).Format(time.RFC3339)
	path := fmt.Sprintf("expiration=%s&id=%s&name=%s", theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(hmacSecret))
	_, err := h.Write([]byte(path))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	symlink := fmt.Sprintf("/symlink/take?hash=%s&%s", hash, path)

	return symlink
}
func Test_TakeSymlink_happy_path(t *testing.T) {
	fileContent := bytes.NewBuffer([]byte("FileContent"))

	md := storage.NewMetadata("Breaker")
	hmacSecret := "MySecret"
	mockSymlink := formatSymlink(md, hmacSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		_, err := io.Copy(file, fileContent)
		require.NoError(t, err)

		return nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)
	handlers.TakeSymlink(fileDownloader, hmacSecret, symlinkParamsx).ServeHTTP(w, r)

	// LOL why cant i copy this?
	fileContent = bytes.NewBuffer([]byte("FileContent"))
	require.Equal(t, fileContent, w.Body)

	require.Equal(t, http.StatusOK, w.Code)

}

func Test_TakeSymlink_unauthentic_symlink(t *testing.T) {

	md := storage.NewMetadata("Breaker")
	dirtyHmacSecret := "DirtySecret"
	mockSymlink := formatSymlink(md, dirtyHmacSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		return nil
	}

	realHmacSecret := "RealSecret"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)
	handlers.TakeSymlink(fileDownloader, realHmacSecret, symlinkParamsx).ServeHTTP(w, r)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func Test_TakeSymlink_expired_symlink(t *testing.T) {

	md := storage.NewMetadata("Breaker")
	hmacSecret := "MySecret"
	mockSymlink := formatSymlinkWithBadTime(md, hmacSecret)

	fileDownloader := func(file io.Writer, filename string) error {
		return nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, mockSymlink, nil)
	handlers.TakeSymlink(fileDownloader, hmacSecret, symlinkParamsx).ServeHTTP(w, r)

	require.Equal(t, http.StatusNotAcceptable, w.Code)
}
