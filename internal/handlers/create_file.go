package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyqe/brigand/internal/storage"
)

func NewCreateFile(
	metadataClient storage.MetadataClient,
	upload storage.FileUploader,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "file content required", http.StatusBadRequest)
			return
		}

		filename := r.URL.Query().Get("name")
		if filename == "" {
			http.Error(w, "filename required", http.StatusBadRequest)
			return
		}

		md := storage.NewMetadata(filename)

		err := metadataClient.Create(r.Context(), md)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = upload(r.Body, md.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(md)
	}
}
