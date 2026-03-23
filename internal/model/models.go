package model

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
type Event struct {
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

type ResponseDataBatchByCookie struct {
	OriginalURL string `json:"original_url"` // "<оригинальны URL>",
	ShortURL    string `json:"short_url"`    // "<результирующий сокращённый URL>"
}
