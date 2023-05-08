package handlers

import "net/http"

func NewGetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("I was here")
		w.WriteHeader(http.StatusOK)
	}
}
