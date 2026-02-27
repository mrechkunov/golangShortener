package repository_test

import (
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestSafeSlice_SetData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url      string
		shortURL string
	}{
		{"Valid name", "http://test.test1", "http://test.short"},
		{"Empty name", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.Storage.SetData(tt.url, tt.shortURL)
		})
	}
}

func TestSafeSlice_GetData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		shortURL string
		want     string
	}{
		{"Valid name", "http://test.test2", "http://test.test2"},
		{"Empty name", "", ""},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.Storage.SetData(tt.want, tt.shortURL)
			got, _ := repository.Storage.GetData(tt.shortURL)

			// TODO: update the condition below to compare got with tt.want.
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("SelectData() = %v want %v", got, tt.want)
			}
		})
	}
}
