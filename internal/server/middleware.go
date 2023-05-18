package server

import (
	"net/http"
)

const (
	Byte     int64 = 1
	Kilobyte       = Byte * 1000
	Megabyte       = Kilobyte * 1000
	Gigabyte       = Megabyte * 1000
)

// MaxUpload
func MaxUpload(max int64) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, max)
			next.ServeHTTP(w, r)
		}
	}
}
