package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"go.uber.org/zap"
)

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()

	config.Init()
	//log.Println("reading config")
	r := chi.NewRouter()
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	Sugar.Infow(
		"Starting server",
		"addr", config.ConfigAdreses.ServerBindAdress,
	)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		Sugar.Fatalw(err.Error(), "event", "start server")
	}
}
