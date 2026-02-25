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
	E  []Event
}

func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		E: make([]Event, 0),
	}
}

func (s *SafeSlice) SetData(shortURL string, url string) {
	p, err := NewProducer(config.ConfigAdreses.JSONFile) // создаем новый продюсер для записи в файл
	if err != nil {
		logger.Log.Errorln("error while file opening (Producer)")
	}
	defer p.Close() // закроем файл при выходе из функции

	s.mu.Lock()         // блокировка на запись
	defer s.mu.Unlock() // разблокировка при выходе из функции

	var nextElement Event
	if len(s.E) == 0 {
		nextElement.ID = 1
	} else {
		nextElement.ID = s.E[len(s.E)-1].ID + 1
	}
	nextElement.OriginalURL = url
	nextElement.ShortURL = shortURL
	// проверяем по url есть ли в слайсе такой элемент
	isExist := false
	for _, el := range s.E {
		if el.OriginalURL == nextElement.OriginalURL {
			isExist = true
		}
	}
	if isExist {
		logger.Log.Infoln("URL", nextElement.OriginalURL, "already exist in storage")
	} else {
		s.E = append(s.E, nextElement)
		p.WriteEvent(&nextElement)
	}
}

func (s *SafeSlice) ReadDataFromFile() {
	c, err := NewConsumer(config.ConfigAdreses.JSONFile)
	if err != nil {
		logger.Log.Errorln("error while file opening (Consumer)")
	}
	defer c.Close()

	for {
		el, err := c.ReadEvent()
		if err != nil {
			break
		} else {
			s.mu.Lock() // блокировка на запись

			var nextElement Event
			if len(s.E) == 0 {
				nextElement.ID = 1
			} else {
				nextElement.ID = s.E[len(s.E)-1].ID + 1
			}

			nextElement.OriginalURL = el.OriginalURL
			nextElement.ShortURL = el.ShortURL
			s.E = append(s.E, nextElement)
			s.mu.Unlock() // разблокировка
		}

		//		fmt.Println(el)
	}
}

// for {
// 	if el, err = c.ReadEvent(); err == nil && el != nil {
// 		fmt.Println(el)

// 		s.mu.Lock() // Блокировка на запись
// 		defer s.mu.Unlock()
// 		s.E = append(s.E, *el)
// 	} else {
// 		break
// 	}
// }

func (s *SafeSlice) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	var urlToReturn string
	isExist := false
	for _, el := range s.E {
		if el.ShortURL == shortURL {
			isExist = true
			urlToReturn = el.OriginalURL
		}
	}
	if !isExist {
		logger.Log.Infoln("ShortURL", shortURL, "is not exist in storage")
	}
	return urlToReturn, isExist
}

var Storage *SafeSlice = NewSafeSlice()
