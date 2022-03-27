package storage

import (
	"context"
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/database"
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var UniqueViolation = DatabaseError{Code: pgerrcode.UniqueViolation}
var NoRowFound = DatabaseError{Err: errors.New("no rows in result set")}

type PostgresRepository struct {
	handler basedbhandler.DBHandler
}

func NewPostgresRepository(handler basedbhandler.DBHandler) (*PostgresRepository, error) {
	var repo PostgresRepository
	repo.handler = handler
	return &repo, nil
}

func (r *PostgresRepository) Ping(ctx context.Context) (bool, error) {
	row, err := r.handler.QueryRow(ctx, "select 10")
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

func (r *PostgresRepository) FindByUser(ctx context.Context, userID string) ([]UserURLs, error) {
	rows, err := r.handler.Query(ctx, database.GetURLsByUserID, userID)
	if err != nil {
		return nil, err
	}
	var resArr []UserURLs
	for rows.Next() {
		var rec UserURLs
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.OriginalURL, &rec.ShortURL)
		resArr = append(resArr, rec)
		if err != nil {
			return nil, err
		}
	}
	return resArr, nil
}

func (r *PostgresRepository) Save(ctx context.Context, userID string, originalURL string, shortURL string) error {
	err := r.handler.Execute(ctx, database.InsertURL, userID, nil, originalURL, shortURL)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &UniqueViolation
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) SaveBatch(ctx context.Context, src UserBatchURLs) error {
	var paramArr [][]interface{}
	for _, obj := range src.List {
		var paramLine []interface{}
		paramLine = append(paramLine, src.UserID)
		paramLine = append(paramLine, obj.CorrelationID)
		paramLine = append(paramLine, obj.OriginalURL)
		paramLine = append(paramLine, obj.ShortURL)
		paramArr = append(paramArr, paramLine)
	}
	err := r.handler.ExecuteBatch(ctx, database.InsertURL, paramArr)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &UniqueViolation
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) FindByShort(ctx context.Context, userID string, shortURL string) (string, error) {
	var err error
	var row basedbhandler.Row
	if userID == "" {
		row, err = r.handler.QueryRow(ctx, database.GetOriginalURLByShort, shortURL)
	} else {
		row, err = r.handler.QueryRow(ctx, database.GetOriginalURLByShortForUser, userID, shortURL)
	}
	if err != nil {
		return "", err
	}
	var res string

	err = row.Scan(&res)
	if err != nil && err.Error() == "no rows in result set" {
		return "", &NoRowFound
	}
	if err != nil {
		return "", err
	}
	return res, nil
}
