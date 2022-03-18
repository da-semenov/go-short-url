package server

import (
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

func (s *UserService) GetURLsByUser(userID string) ([]urls.UserURLs, error) {
	_, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	return []urls.UserURLs{}, nil
}

func (s *UserService) Ping() bool {
	res, _ := s.repository.Ping()
	return res
}
