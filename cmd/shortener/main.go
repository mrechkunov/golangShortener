package main

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/mrechkunov/golangShortener.git/internal/service"
)

func main() {
	config.Init()
	Storage := repository.StorageInit()
	defer Storage.Close()
	defer logger.Log.Sync() // закрываем логгер при выходе из main
	logger.Log.Infoln("Reading config")

	r := chi.NewRouter()

	// mount routes to profiling
	r.Route("/debug", func(r chi.Router) {
		r.Get("/pprof/", pprof.Index)
		r.Get("/pprof/cmdline", pprof.Cmdline)
		r.Get("/pprof/profile", pprof.Profile)
		r.Get("/pprof/symbol", pprof.Symbol)
		r.Get("/pprof/trace", pprof.Trace)
		r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/pprof/heap", pprof.Handler("heap"))
		r.Handle("/pprof/mutex", pprof.Handler("mutex"))
		r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/pprof/block", pprof.Handler("block"))
	})

	// GET Handlers
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handler.GetHandler)))
	r.Get("/ping", logger.WithLogging(gzipMiddleware(handler.GetHandlerPingDB)))
	r.Get("/api/user/urls", logger.WithLogging(gzipMiddleware(handler.GetHandlerURLs)))
	chanToDelete := make(chan []string)
	go service.SetIsDeleted(chanToDelete)
	// DELETE Handlers
	r.Delete("/api/user/urls", logger.WithLogging(gzipMiddleware(handler.DeleteHandler(chanToDelete))))
	// POST Handlers
	r.Post("/", logger.WithLogging(gzipMiddleware(handler.PostHandler)))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handler.JSONPostHandler)))
	r.Post("/api/shorten/batch", logger.WithLogging(gzipMiddleware(handler.JSONBatchPostHandler)))

	logger.Log.Infoln("Starting server", "addr", config.ConfigAdreses.ServerBindAdress)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}
	os.Exit(21)
	close(chanToDelete)
}
