package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/handler"
	"github.com/stretchr/testify/assert"
)

type RequestBody struct {
	URL string `json:"url"`
}

type ResponseBody struct {
	Result string `json:"result"`
}

func TestJSONPostHandler(t *testing.T) {
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
				contentType: "application/json",
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
			requestBody := RequestBody{tt.reqBody}
			body, err := json.Marshal(requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(tt.method, "/", bytes.NewBuffer(body))
			//создаем новый Recorder
			w := httptest.NewRecorder()
			handler.JSONPostHandler(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			if status := w.Code; status != http.StatusCreated {
				assert.Equal(t, tt.want.code, resp.StatusCode)
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			} else {

				expectedResponse := ResponseBody{Result: tt.want.response}
				var actualResponse ResponseBody
				err = json.NewDecoder(w.Body).Decode(&actualResponse)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.want.code, resp.StatusCode)
				assert.Equal(t, expectedResponse, actualResponse)
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			}

		})
	}
}
