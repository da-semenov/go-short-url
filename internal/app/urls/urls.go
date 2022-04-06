package urls

import "errors"

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type UserBatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchDelete = string

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("no rows in result set")
