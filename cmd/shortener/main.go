package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
)

func main() {
	config.Init()
	r := chi.NewRouter()
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)

	//mux := http.NewServeMux()
	//mux.HandleFunc(`/`, handler.PostHandler)
	//mux.HandleFunc(`/{id}`, handler.GetHandler)
	err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r)
	if err != nil {
		panic(err)
	}

}
