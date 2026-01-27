package repository

var Storage = make(map[string]string)

func InsertData(url string, shortUrl string) {
	Storage[shortUrl] = url
}

func SelectData(shortUrl string) string {
	return Storage[shortUrl]
}
