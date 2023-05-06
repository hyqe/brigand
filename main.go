package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func store_a_file(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if len(b) < 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id": uuid.New(),
	})

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/files", store_a_file)

	addr := fmt.Sprintf(":%v", os.Getenv("PORT"))
	fmt.Println("Running on ", addr)
	http.ListenAndServe(addr, router)
}
