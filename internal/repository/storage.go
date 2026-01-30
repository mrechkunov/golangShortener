package repository

import "sync"

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

var Storage *SafeMap = NewSafeMap()

func SetData(shortURL string, url string) {
	Storage.mu.Lock() // Блокировка на запись
	defer Storage.mu.Unlock()
	Storage.m[shortURL] = url
}

func GetData(shortURL string) (string, bool) {
	// Storage.mu.RLock() // Блокировка на чтение
	// defer Storage.mu.RUnlock()
	url, ok := Storage.m[shortURL]
	return url, ok
}
