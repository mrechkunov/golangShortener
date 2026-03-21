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
	Storage := repository.StorageInit()
	defer Storage.Close()
	defer logger.Log.Sync() // закрываем логгер при выходе из main
	logger.Log.Infoln("Reading config")

	r := chi.NewRouter()
	// GET Handlers
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handler.GetHandler)))
	r.Get("/ping", logger.WithLogging(gzipMiddleware(handler.GetHandlerPingDB)))
	r.Get("/api/user/urls", logger.WithLogging(gzipMiddleware(handler.GetHandlerURLs)))

	// POST Handlers
	r.Post("/", logger.WithLogging(gzipMiddleware(handler.PostHandler)))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handler.JSONPostHandler)))
	r.Post("/api/shorten/batch", logger.WithLogging(gzipMiddleware(handler.JSONBatchPostHandler)))

	logger.Log.Infoln("Starting server", "addr", config.ConfigAdreses.ServerBindAdress)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}
}
