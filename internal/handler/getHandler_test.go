package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string // наименовение теста
		want        want   // ожидаемый результат
		reqEndPoint string // точка входа
	}{
		{
			name:        "Bad request test",
			reqEndPoint: "/badRequest",
			want: want{
				code:        400,
				response:    "short URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		// {
		// 	name:        "positive test",
		// 	reqEndPoint: "/",
		// 	want: want{
		// 		code:        307,
		// 		response:    "short URL not found\n",
		// 		contentType: "text/plain; charset=utf-8",
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.reqEndPoint, nil)
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
