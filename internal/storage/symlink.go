package storage

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
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

func MakeSymlinkSignature(expiration, id, name string) string {
	v := url.Values{}
	v.Add("expiration", expiration)
	v.Add("id", id)
	v.Add("name", name)

	return v.Encode()
}

func MakeSymlinkQuery(expiration, id string) string {
	v := url.Values{}
	v.Add("expiration", expiration)
	v.Add("id", id)

	return v.Encode()
}

type Symlink struct {
	Hash       string
	Expiration string
	Id         string
	Name       string
	Link       string
}
