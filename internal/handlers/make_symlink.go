package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyqe/brigand/internal/storage"
	"net/http"
	"time"
)

func formatSymlink(md *storage.Metadata, symlinkSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(time.Hour).Format(time.RFC3339)
	query := fmt.Sprintf("expiration=%s&id=%s&name=%s", theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(symlinkSecret))
	_, err := h.Write([]byte(query))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	symlink := fmt.Sprintf("https://brigand.hyqe.org/symlink/take?hash=%s&%s", hash, query)

	return symlink
}

func MakeSymlink(
	metadataClient storage.MetadataClient,
	upload storage.FileUploader,
	getFileId func(r *http.Request) string,
	symlinkSecret string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// symlink := storage.SymlinkFromQuery(r)
		if r.Body == nil || r.ContentLength < 1 {
			http.Error(w, "file content required", http.StatusBadRequest)
			return
		}

		filename := getFileId(r)
		if filename == "" {
			http.Error(w, "filename required", http.StatusBadRequest)
			return
		}

		md := storage.NewMetadata(filename)
		err := metadataClient.Create(r.Context(), md)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = upload(r.Body, md.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		symlink := map[string]string{
			"Link": formatSymlink(md, symlinkSecret),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(symlink)

	}
}
