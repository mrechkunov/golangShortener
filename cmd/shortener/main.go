package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func main() {
	config.Init()
	defer logger.Log.Sync() // закрываем логгер при выходе из main
	logger.Log.Infoln("Reading config")
	repository.Storage.ReadDataFromFile()
	r := chi.NewRouter()
	r.Post("/", logger.WithLogging(gzipMiddleware(handler.PostHandler)))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handler.JSONPostHandler)))
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handler.GetHandler)))
	logger.Log.Infoln("Starting server", "addr", config.ConfigAdreses.ServerBindAdress)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}
}
