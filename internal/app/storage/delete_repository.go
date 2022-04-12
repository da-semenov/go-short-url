package storage

import (
	"context"
	"github.com/da-semenov/go-short-url/internal/app/database"
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
)

type DeleteRepository struct {
	handler basedbhandler.DBHandler
}

func NewDeleteRepository(handler basedbhandler.DBHandler) (*DeleteRepository, error) {
	var repo DeleteRepository
	repo.handler = handler
	return &repo, nil
}

func (r *DeleteRepository) BatchDelete(ctx context.Context, userID string, URLList []string) error {
	var paramArr [][]interface{}
	for _, l := range URLList {
		var paramLine []interface{}
		if l != "" {
			paramLine = append(paramLine, userID)
			paramLine = append(paramLine, l)
			paramArr = append(paramArr, paramLine)
		}
	}
	err := r.handler.ExecuteBatch(ctx, database.DeleteUserURL, paramArr)
	return err
}
