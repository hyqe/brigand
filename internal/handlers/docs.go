package handlers

import (
	"io"
	"net/http"
	"os"
)

func NewGetDocs(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(filename)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.Copy(w, f)
		f.Close()
	}
}
