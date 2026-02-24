package repository

import (
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"go.uber.org/zap"
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
	// создаём предустановленный регистратор zap
	logg, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logg.Sync()
	// делаем регистратор SugaredLogger
	logger.Sugar = *logg.Sugar()

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
		logger.Sugar.Infoln("URL", nextElement.OriginalURL, "already exist in storage")
		//fmt.Println("URL", nextElement.OriginalURL, "already exist in storage")
	} else {
		s.e = append(s.e, nextElement)
		Producer.WriteEvent(&nextElement)
	}
}

func (s *SafeSlice) ReadDataFromFile() {
	// создаём предустановленный регистратор zap
	logg, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logg.Sync()
	// делаем регистратор SugaredLogger
	logger.Sugar = *logg.Sugar()
	var C *Consumer
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

	// создаём предустановленный регистратор zap
	logg, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logg.Sync()
	// делаем регистратор SugaredLogger
	logger.Sugar = *logg.Sugar()
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
		logger.Sugar.Infoln("ShortURL", shortURL, "is not exist in storage")
		//fmt.Println("ShortURL", shortURL, "is not exist in storage")
	}
	return urlToReturn, isExist
}

var Storage *SafeSlice = NewSafeSlice()
