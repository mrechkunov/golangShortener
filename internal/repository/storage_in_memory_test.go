package repository_test

import (
	"testing"

	"github.com/mrechkunov/golangShortener.git/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestSafeSlice_SetData(t *testing.T) {
	testing.Init()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url      string
		shortURL string
		cookie   string
	}{
		{"Valid name", "http://test.test1", "http://test.short", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851"},
		{"Empty name", "", "", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Storage := repository.StorageInit()
			defer Storage.Close()
			repository.GetStorage().SetData(tt.url, tt.shortURL, tt.cookie)
		})
	}
}

func TestSafeSlice_GetData(t *testing.T) {
	testing.Init()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		shortURL string
		want     string
		cookie   string
	}{
		{"Valid name", "http://test.test2", "http://test.test2", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851"},
		{"Empty name", "", "", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Storage := repository.StorageInit()
			defer Storage.Close()
			repository.GetStorage().SetData(tt.want, tt.shortURL, tt.cookie)
			got, _ := repository.GetStorage().GetData(tt.shortURL)

			// TODO: update the condition below to compare got with tt.want.
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("SelectData() = %v want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMap_IsCookieExist(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cookie string
		want   bool
	}{
		{"notFoundTest", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851", false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := repository.NewSafeMap()
			got := s.IsCookieExist(tt.cookie)
			// TODO: update the condition below to compare got with tt.want.
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("IsCookieExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
