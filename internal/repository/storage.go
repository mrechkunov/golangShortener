package repository

import "sync"

type Event struct {
	Id          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type SafeSlice struct {
	mu sync.RWMutex
	e  []Event
}

func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		e: make([]Event, 0),
	}
}

func (s *SafeSlice) SetData(shortURL string, url string) {
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	if s.e[len(s.e)] == 0 {
		nextID := 1
	} else {
		nextID := s.e[len(s.e)-1].Id - 1
	}

	nextEvent := s[len(e)-1]
	nextEvent[0].Id

	s = append(s)
	s.s[shortURL] = url
}

func (s *SafeMap) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	url, ok := s.s[shortURL]
	return url, ok
}

var Storage *SafeMap = NewSafeMap()
