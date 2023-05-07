package handlers

import (
	"net/http"

	"github.com/hyqe/brigand/internal/storage"
)

func NewGetFileById(
	metadataClient storage.MetadataClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: https://github.com/hyqe/brigand/issues/3
	}
}
