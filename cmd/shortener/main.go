package main

import (
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/handler"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.PostHandler)
	mux.HandleFunc(`/{id}`, handler.GetHandler)
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}

}
