package server

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/models"
	"github.com/da-semenov/go-short-url/internal/app/urls"
)

type EncodeFunc func(str string) string

type UserService struct {
	dbRepository   models.DBRepository
	fileRepository models.FileRepository
	encode         EncodeFunc
	baseURL        string
}

func NewUserService(repoDB models.DBRepository, repoFile models.FileRepository, baseURL string) *UserService {
	var s UserService
	s.dbRepository = repoDB
	s.fileRepository = repoFile
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	s.baseURL = baseURL
	return &s
}

func (s *UserService) GetID(url string) (string, string, error) {
	if url == "" {
		return "", "", errors.New("url is empty")
	}
	key := s.encode(url)
	return s.baseURL + key, key, nil
}

func (s *UserService) mapUserURLs(src *models.UserURLs) (*urls.UserURLs, error) {
	return &urls.UserURLs{ShortURL: s.baseURL + src.ShortURL, OriginalURL: src.OriginalURL}, nil
}

func (s *UserService) GetURLsByUser(ctx context.Context, userID string) ([]urls.UserURLs, error) {
	if userID == "" {
		return nil, errors.New("user_id is empty")
	}
	resArr, err := s.dbRepository.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resList := make([]urls.UserURLs, 0)
	for _, rec := range resArr {
		u, err := s.mapUserURLs(&rec)
		if err != nil {
			return nil, errors.New("can't map result to UserURLs")
		}
		resList = append(resList, *u)
	}
	return resList, nil
}

func (s *UserService) SaveUserURL(ctx context.Context, userID string, originalURL string, shortURL string) error {
	err := s.fileRepository.Save(shortURL, originalURL)
	if err != nil {
		return err
	}

	err = s.dbRepository.Save(ctx, userID, originalURL, shortURL)
	if errors.Is(err, &models.UniqueViolation) {
		return urls.ErrDuplicateKey
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) SaveBatch(ctx context.Context, userID string, src []urls.UserBatch) ([]urls.UserBatchResult, error) {
	var (
		res     models.UserBatchURLs
		err     error
		resurls []urls.UserBatchResult
	)
	res.UserID = userID
	for _, obj := range src {
		var e models.Element
		var fullShortURL string
		e.CorrelationID = obj.CorrelationID
		e.OriginalURL = obj.OriginalURL
		fullShortURL, e.ShortURL, err = s.GetID(obj.OriginalURL)
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, e)
		resurls = append(resurls, urls.UserBatchResult{CorrelationID: obj.CorrelationID, ShortURL: fullShortURL})
	}
	err = s.dbRepository.SaveBatch(ctx, res)
	if errors.Is(err, &models.UniqueViolation) {
		return nil, urls.ErrDuplicateKey
	}
	if err != nil {
		return nil, err
	}
	return resurls, nil
}

func (s *UserService) GetURLByShort(ctx context.Context, userID string, shortURL string) (string, error) {
	if shortURL == "" {
		return "", errors.New("shortURL is empty")
	}
	originalURL, err := s.dbRepository.FindByShort(ctx, userID, shortURL)
	if errors.Is(err, &models.NoRowFound) {
		return "", urls.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s *UserService) Ping(ctx context.Context) bool {
	res, _ := s.dbRepository.Ping(ctx)
	return res
}
