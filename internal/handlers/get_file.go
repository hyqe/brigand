package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyqe/brigand/internal/storage"
)

func NewGetFileById(
	metadataClient storage.MetadataClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: https://github.com/hyqe/brigand/issues/3
		fileId := mux.Vars(r)["fileId"]

		ctx := context.Background()
		md, err := metadataClient.GetById(ctx, fileId)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		magicdata, err := metadataClient.GetFileById(ctx, md.FileName)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "multipart/form-data")
		json.NewEncoder(w).Encode(magicdata)
	}
}
