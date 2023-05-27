package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hyqe/brigand/internal/storage"
)

func formatSymlink(md *storage.Metadata, symlinkSecret string) string {
	// RFC3339 2006-01-02T15:04:05Z07:00
	theTime := time.Now().Add(time.Hour).Format(time.RFC3339)

	signature := storage.MakeSymlinkSignature(theTime, md.Id, md.FileName)

	h := hmac.New(sha256.New, []byte(symlinkSecret))
	_, err := h.Write([]byte(signature))
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(h.Sum(nil))

	// TODO: what is this
	query := storage.MakeSymlinkQuery(theTime, md.Id)
	symlink := fmt.Sprintf("https://brigand.hyqe.org/symlink/take/%s?hash=%s&%s", md.FileName, hash, query)

	return symlink
}

func MakeSymlink(
	metadataClient storage.MetadataClient,
	upload storage.FileUploader,
	getFilename func(r *http.Request) string,
	symlinkSecret string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body == nil || r.ContentLength < 1 {
			http.Error(w, "file content required", http.StatusBadRequest)
			return
		}

		md := storage.NewMetadata(getFilename(r))
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
