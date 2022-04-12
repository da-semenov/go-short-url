package app

import (
	"encoding/base64"
)

type EncodeFunc func(str string) string

type URLService struct {
	repository Repository
	encode     EncodeFunc
	baseURL    string
}

type Repository interface {
	Find(id string) (string, error)
	Save(id string, value string) error
}

func NewService(repo Repository, baseURL string) *URLService {
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

func (s *URLService) GetShorten(url string) (*ShortenResponse, error) {
	var res ShortenResponse
	resStr, err := s.GetID(url)
	if err != nil {
		return nil, err
	}
	res.Result = resStr
	return &res, nil
}
