package repository

import (
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type Event struct {
	ID          int    `json:"uuid"`
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
	Producer, err := NewProducer(config.ConfigAdreses.JSONFile)
	if err != nil {
		logger.Sugar.Errorln("error while file opening (Producer)")
	}
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	var nextElement Event
	if len(s.e) == 0 {
		nextElement.ID = 1
	} else {
		nextElement.ID = s.e[len(s.e)-1].ID + 1
	}
	nextElement.OriginalURL = url
	nextElement.ShortURL = shortURL
	// проверяем по url есть ли в слайсе такой элемент
	isExist := false
	for _, el := range s.e {
		if el.OriginalURL == nextElement.OriginalURL {
			isExist = true
		}
	}
	if isExist {
		logger.Sugar.Infow("URL", nextElement.OriginalURL, "already exist in storage")
	} else {
		s.e = append(s.e, nextElement)
		Producer.WriteEvent(&nextElement)
	}
}

func (s *SafeSlice) ReadDataFromFile() {
	var C *Consumer
	var err error
	C, err = NewConsumer(config.ConfigAdreses.JSONFile)
	if err != nil {
		logger.Sugar.Errorln("error while file opening (Consumer)")
	}
	var el *Event
	for {
		if el, err = C.ReadEvent(); err != nil && el != nil {
			s.mu.Lock() // Блокировка на запись
			defer s.mu.Unlock()
			s.e = append(s.e, *el)
		} else {
			break
		}

	}
}

func (s *SafeSlice) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	var urlToReturn string
	isExist := false
	for _, el := range s.e {
		if el.ShortURL == shortURL {
			isExist = true
			urlToReturn = el.OriginalURL
		}
	}
	if !isExist {
		logger.Sugar.Infow("ShortURL", shortURL, "is not exist in storage")
	}
	return urlToReturn, isExist
}

var Storage *SafeSlice = NewSafeSlice()
