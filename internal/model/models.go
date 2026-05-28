package model

import "time"

type RequestData struct {
	URL string `json:"url"`
}

type ResponseData struct {
	ShortURL string `json:"result"`
}

type RequestBody struct {
	URL string `json:"url"`
}

type ResponseBody struct {
	Result string `json:"result"`
}

//generate:reset
type Event struct { //resetable struct
	ID          int    `json:"event_id" db:"count"`
	ShortURL    string `json:"short_url" db:"shorturl"`
	OriginalURL string `json:"original_url" db:"originalurl"`
	Cookie      string `json:"cookie" db:"cookie"`
	UID         uint32 `json:"user_id" db:"uuid"`
	IsDeleted   bool   `json:"is_deleted" db:"isdeleted"`
}

type RequestDataBatch struct {
	CorrelationID string `json:"correlation_id"` // "<строковый идентификатор>",
	OriginalURL   string `json:"original_url"`   // "<URL для сокращения>"
}

type ResponseDataBatch struct {
	CorrelationID string `json:"correlation_id"` // "<строковый идентификатор из объекта запроса>",
	ShortURL      string `json:"short_url"`      // "<результирующий сокращённый URL>"
}

//generate:reset
type ResponseDataBatchByCookie struct {
	OriginalURL string `json:"original_url"` // "<оригинальны URL>",
	ShortURL    string `json:"short_url"`    // "<результирующий сокращённый URL>"
}

//generate:reset
type ObserverEvent struct {
	Ts          time.Time `json:"ts"`      // : 12345678  unix timestamp события
	Action      string    `json:"action"`  // : "shorten",   // действие: shorten (создание) или follow (прохождение по ссылке)
	UserId      uint32    `json:"user_id"` // : "12315134", // идентификатор пользователя, если есть
	OriginalURL string    `json:"url"`     // : "https://mylongdomain.com/my/long/path/to/shorten/" // оригинальный (не сокращенный) URL
}

// интерфейс пулов
type ResetableStruct interface {
	Reset()
}

// структура пулла
type Pool[T ResetableStruct] struct {
	pool chan T
	new  func() T
}

// NewPool создает новый пул.
// необходимо передать конструктор (newFunc) для создания новых объектов.
func NewPool[T ResetableStruct](capacity int, newFunc func() T) *Pool[T] {
	return &Pool[T]{
		pool: make(chan T, capacity),
		new:  newFunc,
	}
}

// Get получает объект из пула или создает новый, если пул пуст.
func (p *Pool[T]) Get() T {
	select {
	case obj := <-p.pool:
		return obj
	default:
		return p.new()
	}
}

// Put возвращает объект в пул, автоматически вызывая метод Reset().
func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	select {
	case p.pool <- obj:
	default:
		// Пул переполнен, объект отбрасывается и уходит сборщику мусора
	}
}
