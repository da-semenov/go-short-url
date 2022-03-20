package storage

import (
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/database"
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
)

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
