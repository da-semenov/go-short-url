package handlers

import (
	"context"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (s *UserServiceMock) GetURLsByUser(ctx context.Context, userID string) ([]urls.UserURLs, error) {
	args := s.Called(userID)
	var res []urls.UserURLs
	res = append(res, urls.UserURLs{OriginalURL: args.String(0), ShortURL: args.String(0)})
	return res, args.Error(1)
}

func (s *UserServiceMock) Ping(ctx context.Context) bool {
	args := s.Called()
	return args.Bool(0)
}

func (s *UserServiceMock) SaveUserURL(ctx context.Context, userID string, originalURL string, shortURL string) error {
	args := s.Called(userID, originalURL, shortURL)
	if originalURL == "bad_URL" {
		return args.Error(0)
	} else {
		return nil
	}
}

func (s *UserServiceMock) SaveBatch(ctx context.Context, userID string, src []urls.UserBatch) ([]urls.UserBatchResult, error) {
	args := s.Called(userID, src)
	var res []urls.UserBatchResult
	res = append(res, urls.UserBatchResult{CorrelationID: args.String(0), ShortURL: args.String(1)})
	return res, args.Error(2)
}

func (s *UserServiceMock) GetURLByShort(ctx context.Context, shortURL string) (string, error) {
	args := s.Called(shortURL)
	return args.String(0), args.Error(1)
}

func (s *UserServiceMock) GetID(url string) (string, string, error) {
	args := s.Called(url)
	return args.String(0), args.String(1), args.Error(2)
}

type CryptoServiceMock struct {
	mock.Mock
}

func (s *CryptoServiceMock) Validate(token string) (bool, string) {
	args := s.Called(token)
	return args.Bool(0), token
}

func (s *CryptoServiceMock) GetNewUserToken() (string, string, error) {
	args := s.Called()
	return args.String(0), args.String(1), args.Error(2)
}
