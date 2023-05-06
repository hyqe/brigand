package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_b_happy_path(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files", bytes.NewBuffer([]byte("string")))
	store_a_file(w, r)

	// check status
	require.Equal(t, http.StatusOK, w.Code)

	// check id response
	mappy := make(map[string]any)
	json.NewDecoder(w.Body).Decode(&mappy)
	require.NotEmpty(t, mappy["id"])
}

func Test_store_a_file_no_file(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/files", nil)
	store_a_file(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
