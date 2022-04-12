package server

import (
	"context"
	"github.com/da-semenov/go-short-url/internal/app/models"
	"math"
	"sync"
)

type pool struct {
	sync.Mutex
	maxSize     int
	currentSize int
}

func newPool(maxSize int) *pool {
	var p pool
	p.maxSize = maxSize
	return &p
}

func (p *pool) Inc() {
	p.Lock()
	defer p.Unlock()
	p.currentSize++
}

func (p *pool) Dec() {
	p.Lock()
	defer p.Unlock()
	p.currentSize--
}

func (p *pool) LessMax() bool {
	p.Lock()
	defer p.Unlock()
	return p.maxSize > p.currentSize
}

type deleteJob struct {
	UserID string
	part   []string
}

type DeleteService struct {
	pool         *pool
	taskSize     int
	dbRepository models.DeleteRepository
	jobChanel    chan deleteJob
}

func NewDeleteService(repoDB models.DeleteRepository, poolSize int, taskSize int) *DeleteService {
	var s DeleteService
	s.taskSize = taskSize
	s.pool = newPool(poolSize)
	s.dbRepository = repoDB
	s.jobChanel = make(chan deleteJob, poolSize)
	s.startWorkerPool()
	return &s
}

func split(batchSize int, src []string, resCh chan []string) {
	if batchSize <= 0 || len(src) == 0 {
		close(resCh)
		return
	}
	start := 0
	end := int(math.Min(float64(batchSize), float64(len(src)-start)))
	for start <= len(src) {
		resCh <- src[start:end]
		start = end + 1
		end = start + int(math.Min(float64(batchSize), float64(len(src)-start)))
	}
	close(resCh)
}

func (s *DeleteService) DeleteBatch(ctx context.Context, userID string, URLList []string) error {
	chanel := make(chan []string)
	go split(s.taskSize, URLList, chanel)
	for part := range chanel {
		s.jobChanel <- deleteJob{userID, part}
	}
	s.startWorkerPool()
	return nil
}

func (s *DeleteService) startWorkerPool() {
	for s.pool.LessMax() {
		s.pool.Inc()
		go func() {
			defer s.pool.Dec()
			for {
				for job := range s.jobChanel {
					err := s.dbRepository.BatchDelete(context.Background(), job.UserID, job.part)
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
}
