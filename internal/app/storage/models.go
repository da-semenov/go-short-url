package storage

type StoreRecord struct {
	Key   string
	Value string
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type UserURLs struct {
	id          int
	userID      string
	shortURL    string
	originalURL string
}
