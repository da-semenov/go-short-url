package app

type Service interface {
	ShortURL(url string) (string, error)
	LongURL(key string) (string, error)
}
