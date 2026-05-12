package main

import (
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	//"net/http/pprof"
	_ "net/http/pprof"

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
	// pprof Handlers

	// // Подключение pprof маршрутов
	// r.Route("/debug/pprof", func(r chi.Router) {
	// 	r.HandleFunc("/", pprof.Index)
	// 	r.HandleFunc("/cmdline", pprof.Cmdline)
	// 	r.HandleFunc("/profile", pprof.Profile)
	// 	r.HandleFunc("/symbol", pprof.Symbol)
	// 	r.HandleFunc("/trace", pprof.Trace)
	// 	r.HandleFunc("/heap", pprof.Index)
	// })

	// создаём файл журнала профилирования памяти
	var err error
	config.Fmem, err = os.Create(`./profiles/result.pprof`)
	if err != nil {
		panic(err)
	}
	defer config.Fmem.Close()
	runtime.GC() // получаем статистику по использованию памяти
	if err := pprof.WriteHeapProfile(config.Fmem); err != nil {
		panic(err)
	}

	logger.Log.Infoln("Starting server", "addr", config.ConfigAdreses.ServerBindAdress)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}

	close(chanToDelete)
}
