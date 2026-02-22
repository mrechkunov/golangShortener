package main

import (
	"net/http"
	"strings"

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
	r.Post("/", logger.WithLogging(gzipMiddleware(handler.PostHandler)))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handler.JSONPostHandler)))
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handler.GetHandler)))
	logger.Sugar.Infow(
		"Starting server",
		"addr", config.ConfigAdreses.ServerBindAdress,
	)
	if err := http.ListenAndServe(config.ConfigAdreses.ServerBindAdress, r); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		logger.Sugar.Fatalw(err.Error(), "event", "start server")
	}
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		// проверяем что данные application/json или text/html
		dataFormat := r.Header.Get("Content-type")
		jsonDataFormat := strings.Contains(dataFormat, "application/json")
		textHTMLDataFormat := strings.Contains(dataFormat, "text/html")

		if supportsGzip && (jsonDataFormat || textHTMLDataFormat) {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	}
}
