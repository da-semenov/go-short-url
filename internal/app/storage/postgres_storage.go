package storage

import (
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
)

type PostgresRepository struct {
	handler basedbhandler.DBHandler
}

func NewPostgresRepository(handler basedbhandler.DBHandler) (*PostgresRepository, error) {
	var repo PostgresRepository
	repo.handler = handler
	return &repo, nil
}

func (r *PostgresRepository) Ping() (bool, error) {
	row, err := r.handler.QueryRow("select 10")
	if err != nil {
		return false, err
	}
	var res int
	err = row.Scan(&res)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PostgresRepository) FindByUser(key string) ([]UserURLs, error) {
	return nil, nil
}
