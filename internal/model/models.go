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
}
