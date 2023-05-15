package handlers

import (
	"net/http"

	"context"
	"github.com/gorilla/mux"
	"github.com/hyqe/brigand/internal/storage"
)

func getVar(r *http.Request, var_name string) string {
	return mux.Vars(r)[var_name]
}

func NewGetFileById(
	metadataClient storage.MetadataClient,
	fileDownloader storage.FileDownloader,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: https://github.com/hyqe/brigand/issues/3
		fileId := getVar(r, "fileId")

		ctx := context.Background()
		md, err := metadataClient.GetById(ctx, fileId)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		err = fileDownloader(w, md.Id)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		w.Header().Set("name", md.FileName)
	}
}