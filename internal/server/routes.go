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
	fileDownloader storage.FileDownloader,

) http.Handler {
	r := mux.NewRouter()

	// health check for reverse proxy
	r.HandleFunc("/", handlers.NewGetHealth())

	// get docs html file. the "openapi.html" file is generated by the build process.
	// to generate it manually nodejs is required. run "./scripts/compile_openapi.sh"
	r.HandleFunc("/docs", handlers.NewGetDocs("openapi.html")).
		Methods(http.MethodGet)

	// store a file
	r.HandleFunc("/files", handlers.NewCreateFile(metadataClient)).
		Methods(http.MethodPost)

	// get a file by its Id.
	r.HandleFunc("/files/{fileId}", handlers.NewGetFileById(metadataClient, fileDownloader, getFileId)).
		Methods(http.MethodGet).
		Queries()

	return r
}

func getFileId(r *http.Request) string {
	return mux.Vars(r)["fileId"]
}
