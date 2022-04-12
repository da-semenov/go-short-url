package app

import (
	"encoding/base64"
)

type EncodeFunc func(str string) string

type URLService struct {
	repository Repository
	encode     EncodeFunc
}

type Repository interface {
	Find(id string) (string, error)
	Save(id string, value string) error
}

func NewService(repo Repository) *URLService {
	var s URLService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return &s
}

func (s *URLService) GetID(url string) (string, error) {
	const baseURL string = "http://localhost:8080/"
	id := s.encode(url)
	s.repository.Save(id, url)
	return baseURL + id, nil
}

func (s *URLService) GetURL(id string) (string, error) {
	res, err := s.repository.Find(id)
	return res, err
}
