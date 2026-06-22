package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/mrechkunov/golangShortener.git/internal/service"
	"golang.org/x/crypto/acme/autocert"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	config.Init()
	Storage := repository.StorageInit()
	defer Storage.Close()
	defer logger.Log.Sync() // закрываем логгер при выходе из main
	logger.Log.Infoln("Reading config")
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
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
	// Создаем канал для получения системных сигналов

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// конструируем менеджер TLS-сертификатов
	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist("localhost"),
	}
	var server = &http.Server{
		Addr:      config.ConfigAdreses.ServerBindAdress,
		Handler:   r,
		TLSConfig: manager.TLSConfig(),
	}
	// конструируем сервер
	if config.ConfigAdreses.HttpsEnable {
		logger.Log.Infoln("server starting:", config.ConfigAdreses.ServerBindAdress, "https")
		go func() {
			if err := server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Fatalln(err.Error())
			}

		}()

	} else {
		logger.Log.Infoln("server starting:", config.ConfigAdreses.ServerBindAdress, "http")
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Fatalln(err.Error())
			}

		}()
	}
	// ловим сигналы
	<-stop
	logger.Log.Infoln("Получен сигнал завершения. Начинаем graceful shutdown...")
	// Создаем контекст с таймаутом для завершения активных запросов
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Пытаемся плавно остановить сервер
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Infoln("Сервер завершился с ошибкой:", err)
	} else {
		logger.Log.Infoln("Сервер остановлен корректно.")
	}
	cancel()
	close(chanToDelete)
}
