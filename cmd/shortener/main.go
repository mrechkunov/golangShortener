package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/logger"

	"go.uber.org/zap"
)

func main() {
	// создаём предустановленный регистратор zap
	logg, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logg.Sync()

	// делаем регистратор SugaredLogger
	logger.Sugar = *logg.Sugar()

	config.Init()
	//log.Println("reading config")
	r := chi.NewRouter()
	r.Post("/", logger.WithLogging(handler.PostHandler))
	r.Get("/{id}", logger.WithLogging(handler.GetHandler))
	logger.Sugar.Infow(
		"Starting server",
		"addr", config.ConfigAdreses.ServerBindAdress,
	)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		logger.Sugar.Fatalw(err.Error(), "event", "start server")
	}
}
