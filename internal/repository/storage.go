package repository

import (
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type SafeMap struct {
	mu      sync.RWMutex
	m       map[string]model.Event
	counter int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m:       make(map[string]model.Event),
		counter: 1,
	}
}

func (s *SafeMap) SetData(shortURL string, originalURL string) {
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	newEvent := model.Event{
		ID:          s.counter,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	if _, ok := s.m[shortURL]; !ok {
		s.m[shortURL] = newEvent
		s.counter++
		// запишем событие в файл
		p, err := NewProducer(config.ConfigAdreses.JSONFile) // создаем новый продюсер для записи в файл
		if err != nil {
			logger.Log.Errorln("error while file opening (Producer)")
		}
		defer p.Close() // закроем файл при выходе из функции
		p.WriteEvent(&newEvent)
	} else {
		logger.Log.Infoln("URL", originalURL, "already exist in storage")
	}

}

func (s *SafeMap) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	if val, ok := s.m[shortURL]; !ok {
		return "", false
	} else {
		originalURL := val.OriginalURL
		return originalURL, true
	}
}

var Storage *SafeMap = NewSafeMap()

func (s *SafeMap) ReadDataFromFile() {
	c, err := NewConsumer(config.ConfigAdreses.JSONFile)
	if err != nil {
		logger.Log.Errorln("error while file opening (Consumer)")
	}
	defer c.Close()
	var e *model.Event

	for {
		if e, err = c.ReadEvent(); err != nil {
			logger.Log.Infoln("EOF")
			break
		} else {
			s.mu.Lock() // Блокировка на запись
			s.m[e.ShortURL] = *e
			s.counter = e.ID + 1
			s.mu.Unlock()
		}
	}
}
