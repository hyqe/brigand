package handlers

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed docs.html
var docs string

func GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, docs)
}
