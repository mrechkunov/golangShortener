package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(event *model.Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *Producer) WriteEvents(events *[]model.Event) error {
	data, err := json.MarshalIndent(events, "", "")
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *Producer) Close() {
	p.file.Close()
}

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*model.Event, error) {
	// читаем данные до символа переноса строки
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// преобразуем данные из JSON-представления в структуру
	event := model.Event{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (c *Consumer) ReadEvents() (*[]model.Event, error) {
	// читаем данные до символа переноса строки
	data, err := io.ReadAll(c.reader)
	if err != nil {
		return nil, err
	}

	// преобразуем данные из JSON-представления в структуру
	var events []model.Event
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}
	return &events, nil
}

func (c *Consumer) Close() {
	c.file.Close()
}

type SafeMapFile struct {
	mu      sync.RWMutex
	m       map[string]model.Event
	counter int
}

func NewSafeMapFile() *SafeMapFile {
	var s = SafeMapFile{
		m:       make(map[string]model.Event),
		counter: 1,
	}
	s.ReadDataFromFile()
	return &s
}

func (s *SafeMapFile) SetData(shortURL string, originalURL string, cookie string) error {
	logger.Log.Infoln("Set Data file")
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	uid, err := cryptoauth.GetIDFromCookie(cookie)
	if err != nil {
		logger.Log.Errorln("can not Get ID from cookie while setdata in storage")
	}
	newEvent := model.Event{
		ID:          s.counter,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		Cookie:      cookie,
		UID:         uid,
	}

	if _, ok := s.m[shortURL]; !ok {
		s.m[shortURL] = newEvent
		s.counter++
		// перезапишем файл с новым событием
		p, err := NewProducer(config.ConfigAdreses.JSONFile) // создаем новый продюсер для записи в файл
		if err != nil {
			logger.Log.Errorln("error while file opening (Producer)")
		}
		defer p.Close() // закроем файл при выходе из функции
		var dataSlice []model.Event
		for _, el := range s.m {
			dataSlice = append(dataSlice, el)
		}
		p.WriteEvents(&dataSlice)
	} else {
		logger.Log.Infoln("URL", originalURL, "already exist in storage")
		return errors.New("409 Conflict")
	}
	return nil
}

func (s *SafeMapFile) GetData(shortURL string) (string, bool) {
	logger.Log.Infoln("Get Data file")
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	if val, ok := s.m[shortURL]; !ok {
		return "", false
	} else {
		originalURL := val.OriginalURL
		return originalURL, true
	}
}

func (s *SafeMapFile) ReadDataFromFile() {
	if config.ConfigAdreses.JSONFile == "" {
		logger.Log.Errorln("no file setup (Consumer)")
	} else {
		c, err := NewConsumer(config.ConfigAdreses.JSONFile)
		if err != nil {
			logger.Log.Errorln("error while file opening (Consumer)")
		}
		defer c.Close()

		events, err := c.ReadEvents()
		if err != nil {
			logger.Log.Infow("file is empty (Consumer)")
		} else {
			for _, event := range *events {
				s.mu.Lock() // Блокировка на запись
				s.m[event.ShortURL] = event
				s.mu.Unlock()
			}
		}
	}
}

func (s *SafeMapFile) Close() error {
	return nil
}

func (s *SafeMapFile) IsCookieExist(cookie string) bool {
	// перебор всей мапы и сравнение поля cookie
	returnValue := false
	for _, value := range s.m {
		if value.Cookie == cookie {
			returnValue = true
		}
	}
	return returnValue
}

func (s *SafeMapFile) GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie {
	var result []model.ResponseDataBatchByCookie
	for _, value := range s.m {
		if value.UID == uid {
			var add = model.ResponseDataBatchByCookie{
				OriginalURL: value.OriginalURL,
				ShortURL:    value.ShortURL,
			}
			result = append(result, add)
		}
	}
	return result
}
