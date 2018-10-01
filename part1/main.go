package main // import "github.com/xitonix/xhttp/part1"

import (
	"encoding/json"
	"net/http"

	"github.com/xitonix/xhttp/part1/models"
)

func main() {
	http.Handle("/", &homePageHandler{})
	http.ListenAndServe("", nil)
}

type homePageHandler struct{}

func (*homePageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Language", "go")
	response := models.Response{
		Message: "This is your first HTTP API endpoint in Go",
	}
	w.WriteHeader(http.StatusAccepted)
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
