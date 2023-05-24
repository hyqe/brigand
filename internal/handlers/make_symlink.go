package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"fmt"
	"github.com/hyqe/brigand/internal/storage"
	// "io"
	"net/http"
	"time"
)

// func doesThisHaveFileContent(rc io.ReadCloser) bool {

// 	byties, err := io.ReadAll(rc)
// 	if err != nil {
// 		return false
// 	}

// 	if len(byties) < 1 {
// 		return false
// 	}

// 	return true
// }

func createMetadataInDB(metadataClient storage.MetadataClient, upload storage.FileUploader, getFileName func(r *http.Request) string, w http.ResponseWriter, r *http.Request) (*storage.Metadata, error) {

	//if r.Body == nil || !doesThisHaveFileContent(r.Body) {
	if r.Body == nil || r.ContentLength < 1 {
		http.Error(w, "file content required", http.StatusBadRequest)
		return nil, fmt.Errorf("File content is nothing")
	}

	filename := getFileName(r)
	if filename == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return nil, fmt.Errorf("")
	}

	md := storage.NewMetadata(filename)

	err := metadataClient.Create(r.Context(), md)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, fmt.Errorf("")
	}

	return md, nil
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

	symlink := fmt.Sprintf("https://brigand.hyqe.org/symlink/take?hash=%s&%s", hash, path)

	return symlink
}

type Symlink struct {
	Link string
}

func MakeSymlink(
	metadataClient storage.MetadataClient,
	upload storage.FileUploader,
	getFileId func(r *http.Request) string,
	hmacSecret string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		md, err := createMetadataInDB(metadataClient, upload, getFileId, w, r)
		if err != nil {
			return
		}

		err = upload(r.Body, md.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		symlink := Symlink{
			Link: formatSymlink(md, hmacSecret),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(symlink)

	}
}
