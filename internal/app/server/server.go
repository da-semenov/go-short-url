package server

import (
	"encoding/base64"
	"github.com/da-semenov/go-short-url/internal/app/storage"
)

type EncodeFunc func(str string) string

type URLService struct {
	repository storage.Repository
	encode     EncodeFunc
	baseURL    string
}

func NewService(repo storage.Repository, baseURL string) *URLService {
	var s URLService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	s.baseURL = baseURL
	return &s
}

func (s *URLService) GetID(url string) (string, error) {
	id := s.encode(url)
	err := s.repository.Save(id, url)
	if err != nil {
		return "", err
	}
	return s.baseURL + id, nil
}

func (s *URLService) GetURL(id string) (string, error) {
	res, err := s.repository.Find(id)
	return res, err
}
