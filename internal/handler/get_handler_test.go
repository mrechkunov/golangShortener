package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	testing.Init()
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string // наименовение теста
		want        want   // ожидаемый результат
		reqEndPoint string // точка входа
		body        string // request body
	}{
		{
			name:        "Bad request test",
			reqEndPoint: "/badRequest",
			body:        "ya.ru",
			want: want{
				code:        400,
				response:    "short URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "test not found",
			reqEndPoint: "/",
			body:        "",
			want: want{
				code:        400,
				response:    "short URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Storage := repository.StorageInit()
			defer Storage.Close()
			request := httptest.NewRequest(http.MethodGet, tt.reqEndPoint, strings.NewReader(tt.body))
			//создаем новый Recorder
			w := httptest.NewRecorder()
			handler.GetHandler(w, request)
			res := w.Result()

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			// проверяем
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func BenchmarkGetHandler(b *testing.B) {
	// 1. Создаем фиктивный запрос
	req := httptest.NewRequest("POST", "http://localhost:8080/api/user/urls", nil)

	// 2. Цикл бенчмарка: b.N - количество итераций, которое Go подбирает автоматически
	for i := 0; i < b.N; i++ {
		// 3. Создаем ResponseRecorder для записи ответа
		w := httptest.NewRecorder()

		// 4. Вызываем обработчик напрямую
		handler.GetHandler(w, req)
	}
}
