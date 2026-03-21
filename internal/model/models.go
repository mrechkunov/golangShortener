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
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Cookie      string `json:"cookie"`
	UID         uint32 `json:"user_id"`
}

type RequestDataBatch struct {
	CorrelationID string `json:"correlation_id"` // "<строковый идентификатор>",
	OriginalURL   string `json:"original_url"`   // "<URL для сокращения>"
}

type ResponseDataBatch struct {
	CorrelationID string `json:"correlation_id"` // "<строковый идентификатор из объекта запроса>",
	ShortURL      string `json:"short_url"`      // "<результирующий сокращённый URL>"
}
