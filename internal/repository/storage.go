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

func (s *SafeMap) SetData(shortURL string, url string) {
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	s.m[shortURL] = url
}

func (s *SafeMap) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	url, ok := s.m[shortURL]
	return url, ok
}

var Storage *SafeMap = NewSafeMap()
