package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyqe/brigand/internal/handlers"
	"github.com/hyqe/brigand/internal/storage"
)

// Routes defines all the routes this service offers
func Routes(
	metadataClient storage.MetadataClient,
) http.Handler {
	r := mux.NewRouter()

	// health check for reverse proxy
	r.HandleFunc("/", handlers.NewGetHealth())

	// store a file
	r.HandleFunc("/files", handlers.NewCreateFile(metadataClient)).
		Methods(http.MethodPost)

	// get a file by its Id.
	r.HandleFunc("/files/{fileId}", handlers.NewGetFileById(metadataClient)).
		Methods(http.MethodGet).
		Queries()

	return r
}
