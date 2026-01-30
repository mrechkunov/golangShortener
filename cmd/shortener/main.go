package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
)

func main() {
	config.Init()
	log.Println("reading config")
	r := chi.NewRouter()
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	log.Println("starting web server at:", config.ConfigAdreses.ServerBindAdress)
	err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r)
	if err != nil {
		panic(err)
	}

}
