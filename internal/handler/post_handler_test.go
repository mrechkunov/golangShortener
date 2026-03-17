package handler_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	testing.Init()
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string // наименовение теста
		want    want   // ожидаемый результат
		method  string // метод запроса
		reqBody string // передаваемые параметры запроса
	}{
		{
			name:    "positive test #1 201",
			method:  "POST",
			reqBody: "ya.ru",
			want: want{
				code:        201,
				response:    "http://localhost:8080/7c4e7828",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Bad request #1 400",
			method:  "GET",
			reqBody: "ya.ru",
			want: want{
				code:        400,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Storage := repository.StorageInit()
			defer Storage.Close()
			reqBody := strings.NewReader(tt.reqBody)
			request := httptest.NewRequest(tt.method, "/", reqBody)
			//создаем новый Recorder
			w := httptest.NewRecorder()
			handler.PostHandler(w, request)
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
