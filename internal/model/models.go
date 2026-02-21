package model

type RequestData struct {
	URL string `json:"url"`
}

type ResponseData struct {
	ShortURL string `json:"result"`
}
