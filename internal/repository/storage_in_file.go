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

// Producer to write in file
type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

// NewProducer return ptr to new Producer
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

// WriteEvent write new event to buffer, add endline, write buffer to file
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

// WriteEvents write new events to buffer, add endline, write buffer to file
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

// Close file
func (p *Producer) Close() {
	p.file.Close()
}

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	reader *bufio.Reader
}

// NewConsumer returns new consumer to read from file
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

// ReadEvent read event from file
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

// ReadEvents read events from file
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

// Close Consumer
func (c *Consumer) Close() {
	c.file.Close()
}

type SafeMapFile struct {
	mu      sync.RWMutex
	m       map[string]model.Event
	counter int
}

// NewSafeMapFile return pth to new map with mutex
func NewSafeMapFile() *SafeMapFile {
	var s = SafeMapFile{
		m:       make(map[string]model.Event),
		counter: 1,
	}
	s.ReadDataFromFile()
	return &s
}

// SetData insert to map new row with shortURL, originalURL, UID and rewrite it in to file
func (s *SafeMapFile) SetData(shortURL string, originalURL string, cookie string) error {
	logger.Log.Infoln("Set Data file")
	s.mu.Lock()
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

// Return originalURL from DB by shortURL if row is not exist, return isFound = false
func (s *SafeMapFile) GetData(shortURL string) (originalURL string, isFound bool) {
	logger.Log.Infoln("Get Data file")
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.m[shortURL]; !ok {
		return "", false
	} else {
		originalURL := val.OriginalURL
		return originalURL, true
	}
}

// ReadDataFromFile read all data from file to map
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
				s.mu.Lock()
				s.m[event.ShortURL] = event
				s.mu.Unlock()
			}
		}
	}
}

func (s *SafeMapFile) Close() error {
	return nil
}

// IsCookieExist return true if cookie is exist in map
func (s *SafeMapFile) IsCookieExist(cookie string) bool {
	// перебор всей мапы и сравнение поля cookie
	s.mu.RLock()
	defer s.mu.RUnlock()
	returnValue := false
	for _, value := range s.m {
		if value.Cookie == cookie {
			returnValue = true
		}
	}
	return returnValue
}

// GetDataByUID return batch of URLs whitch user set.
func (s *SafeMapFile) GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie {
	s.mu.RLock()
	defer s.mu.RUnlock()
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

// IsDeleted return true if shortURL is mark as deleted
func (s *SafeMapFile) IsDeleted(shortURL string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[shortURL].IsDeleted
}

// IsCreator return true if user is creator of shortURL else false
func (s *SafeMapFile) IsCreator(shortURL string, cookie string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.m[shortURL].Cookie == cookie {
		return true
	}
	return false
}

// SetIsDeleted  mark all shortURLs from slice as deleted
func (s *SafeMapFile) SetIsDeleted(shortURLs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, val := range shortURLs {
		tmpm := s.m[val]
		tmpm.IsDeleted = true
		s.m[val] = tmpm
	}
}

func (s *SafeMapFile) GetStatData() (statdata model.ResponseStatData) {
	// Создаем map для отслеживания уникальных значений
	uniqueMap := make(map[uint32]struct{})
	s.mu.RLock()
	defer s.mu.RUnlock()
	// перебираем мапу и добавляем только уникальные данные в уникальную мапу
	for _, data := range s.m {
		uniqueMap[data.UID] = struct{}{}
	}
	statdata.Urls = len(s.m)        // количество ключей (коротких URL в базе) так как они уникальны
	statdata.Users = len(uniqueMap) // количество пользователей это размер уникальной мапы
	return
}
