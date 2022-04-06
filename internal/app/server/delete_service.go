package server

import (
	"context"
	"github.com/da-semenov/go-short-url/internal/app/storage"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"sync"
)

type pool struct {
	sync.Mutex
	maxSize     int
	currentSize int
}

type deleteJob struct {
	UserID string
	part   []urls.BatchDelete
}

type DeleteService struct {
	pool         *pool
	dbRepository storage.DeleteRepository
	jobChanel    chan deleteJob
}

func NewDeleteService(repoDB storage.DeleteRepository) *DeleteService {
	var s DeleteService
	s.dbRepository = repoDB
	s.jobChanel = make(chan deleteJob)
	return &s
}

func (s *DeleteService) DeleteBatch(ctx context.Context, userID string, URLList []urls.BatchDelete) error {
	return nil
}
