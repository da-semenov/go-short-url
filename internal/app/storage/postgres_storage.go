package storage

import (
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/database"
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

func (r *PostgresRepository) FindByUser(userID string) ([]UserURLs, error) {
	rows, err := r.handler.Query(database.GetURLsByUserID, userID)
	if err != nil {
		return nil, err
	}
	var resArr []UserURLs
	for rows.Next() {
		var rec UserURLs
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.ShortURL, &rec.OriginalURL)
		resArr = append(resArr, rec)
		if err != nil {
			return nil, err
		}
	}
	return resArr, nil
}

func (r *PostgresRepository) Save(userID string, shortURL string, originalURL string) error {
	err := r.handler.Execute(database.InsertUserURL, userID, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func InitDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("database structure created successfully")
	return nil
}

func ClearDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
