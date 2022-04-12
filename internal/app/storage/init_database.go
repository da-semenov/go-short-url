package storage

import (
	"context"
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/database"
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
)

func InitDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, database.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("database structure created successfully")
	return nil
}

func ClearDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, database.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
