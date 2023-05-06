package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/hyqe/brigand/internal/storage"
)

// CreateFile TODO: add docs here...
func NewCreateFile(
	// TODO: accept db and fs interfaces here.
	metadataClient storage.MetadataClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// - check request content-type: multipart/form-data
		// - store file metadata in db with unique id.
		// - copy file to fs for using unique id as file name.

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		if len(b) < 1 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// TODO: metadataClient.Create(ctx, &storage.Metadata{})

		// TODO:
		// - set response content-type
		// - define CreateFileResponse struct instead of map.
		json.NewEncoder(w).Encode(map[string]any{
			"id": uuid.New(),
		})
	}
}
