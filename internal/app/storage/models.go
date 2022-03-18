package storage

type Repository interface {
	Find(key string) (string, error)
	Save(key string, value string) error
	FindByUser(key string) ([]UserURLs, error)
	Ping() (bool, error)
}

type Repository2 interface {
	FindByUser(key string) ([]UserURLs, error)
	Ping() (bool, error)
}

type UserURLs struct {
	ID          int
	UserID      string
	ShortURL    string
	OriginalURL string
}
