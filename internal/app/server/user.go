package server

import (
	"encoding/base64"
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/storage"
	"github.com/da-semenov/go-short-url/internal/app/urls"
)

type UserService struct {
	repository storage.Repository2
	encode     EncodeFunc
	baseURL    string
}

func NewUserService(repo storage.Repository2, baseURL string) *UserService {
	var s UserService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	s.baseURL = baseURL
	return &s
}

func (s *UserService) GetID(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}
	key := s.encode(url)
	return s.baseURL + key, nil
}

func (s *UserService) mapUserURLs(src *storage.UserURLs) (*urls.UserURLs, error) {
	return &urls.UserURLs{ShortURL: src.ShortURL, OriginalURL: src.OriginalURL}, nil
}

func (s *UserService) mapToUserBatch(src *urls.UserBatch) (*storage.UserBatchURLs, error) {

	return nil, nil
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

func (s *UserService) SaveBatch(userID string, src []urls.UserBatch) ([]urls.UserBatchResult, error) {
	var (
		res     storage.UserBatchURLs
		err     error
		resurls []urls.UserBatchResult
	)
	res.UserID = userID
	for _, obj := range src {
		var e storage.Element
		e.CorrelationID = obj.CorrelationID
		e.OriginalURL = obj.OriginalURL
		e.ShortURL, err = s.GetID(obj.OriginalURL)
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, e)
		resurls = append(resurls, urls.UserBatchResult{CorrelationID: obj.CorrelationID, ShortURL: e.ShortURL})
	}
	err = s.repository.SaveBatch(res)
	if err != nil {
		return nil, err
	}
	return resurls, nil
}

func (s *UserService) Ping() bool {
	res, _ := s.repository.Ping()
	return res
}
