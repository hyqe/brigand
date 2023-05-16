package handlers

import (
	"net/http"

	"context"
	"github.com/hyqe/brigand/internal/storage"
)

func NewGetFileById(
	metadataClient storage.MetadataClient,
	fileDownloader storage.FileDownloader,
	getFileId func(r *http.Request) string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: https://github.com/hyqe/brigand/issues/3
		fileId := getFileId(r)

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