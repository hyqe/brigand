package storage

import (
	"github.com/gorilla/mux"
	"net/http"
	// "net/url"
)

func SymlinkFromQuery(r *http.Request) *Symlink {
	q := r.URL.Query()

	return &Symlink{
		Hash:       q.Get("hash"),
		Expiration: q.Get("expiration"),
		Id:         q.Get("id"),
		Name:       mux.Vars(r)["name"],
	}
}

type Symlink struct {
	Hash       string
	Expiration string
	Id         string
	Name       string
	Link       string
}
