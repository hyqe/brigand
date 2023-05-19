package server_test

import (
	"github.com/hyqe/brigand/internal/server"
	"github.com/stretchr/testify/require"
	"net/http"

	"net/http/httptest"
	"testing"
)

func Test_SudoMiddlware_missing_auth_header(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	mockCredentials := server.Credentials{
		Username: "good_username",
		Password: "good_password",
	}

	mockHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})

	sudo := server.SudoMiddlware(mockCredentials)(mockHandler)
	sudo.ServeHTTP(w, r)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test_SudoMiddlware_deny_request(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetBasicAuth("bad_username", "bad_password")

	mockCredentials := server.Credentials{
		Username: "good_username",
		Password: "good_password",
	}

	mockHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})

	sudo := server.SudoMiddlware(mockCredentials)(mockHandler)
	sudo.ServeHTTP(w, r)

	require.Equal(t, http.StatusUnauthorized, w.Code)

}

func Test_SudoMiddlware_approve_request(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetBasicAuth("good_username", "good_password")

	mockCredentials := server.Credentials{
		Username: "good_username",
		Password: "good_password",
	}

	mockHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})

	sudo := server.SudoMiddlware(mockCredentials)(mockHandler)
	sudo.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)

}
