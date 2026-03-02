package repository

type Store interface {
	GetData(shortURL string) (string, bool)
	SetData(shortURL string, originalURL string)
}
