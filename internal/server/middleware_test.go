package server_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyqe/brigand/internal/server"
	"github.com/stretchr/testify/require"
)

func TestMaxUpload_within_max(t *testing.T) {
	max := server.MaxUpload(server.Byte * 10)

	file := make([]byte, 10)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(file))

	max(mockHandler).ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestMaxUpload_overload(t *testing.T) {
	max := server.MaxUpload(server.Byte * 10)

	file := make([]byte, 11)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.Error(t, err)
		w.WriteHeader(http.StatusBadRequest)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(file))

	max(mockHandler).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
