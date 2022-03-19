package server

import (
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/storage"
	"github.com/da-semenov/go-short-url/internal/app/urls"
)

type UserService struct {
	repository storage.Repository2
}

func NewUserService(repo storage.Repository2) *UserService {
	var s UserService
	s.repository = repo
	return &s
}

func (s *UserService) mapUserURLs(src *storage.UserURLs) (*urls.UserURLs, error) {
	return &urls.UserURLs{ShortURL: src.ShortURL, OriginalURL: src.OriginalURL}, nil
}

func (s *UserService) GetURLsByUser(userID string) ([]urls.UserURLs, error) {
	if userID == "" {
		return nil, errors.New("user_id is empty")
	}
	resArr, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	var resList []urls.UserURLs
	for _, rec := range resArr {
		u, err := s.mapUserURLs(&rec)
		if err != nil {
			return nil, errors.New("can't map result to UserURLs")
		}
		resList = append(resList, *u)
	}
	return resList, nil
}

func (s *UserService) Save(userID string, originalURL string, shortURL string) error {
	err := s.repository.Save(userID, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Ping() bool {
	res, _ := s.repository.Ping()
	return res
}
