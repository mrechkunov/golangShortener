package service

import "github.com/mrechkunov/golangShortener.git/internal/repository"

func SetIsDeleted(chaneltodelete chan []string) {
	for val := range chaneltodelete {
		repository.GetStorage().SetIsDeleted(val)
	}
}
