package storage

import (
	"context"
	"github.com/da-semenov/go-short-url/internal/app/storage/basedbhandler"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresHandler struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

type PostgresRow struct {
	Rows *pgx.Row
}

func (handler *PostgresHandler) Execute(statement string, args ...interface{}) error {
	conn, err := handler.pool.Acquire(handler.ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		_, err = conn.Exec(handler.ctx, statement, args...)
	} else {
		_, err = conn.Exec(handler.ctx, statement)
	}
	return err
}

func (handler *PostgresHandler) QueryRow(statement string, args ...interface{}) (basedbhandler.Row, error) {
	var row pgx.Row
	conn, err := handler.pool.Acquire(handler.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if len(args) > 0 {
		row = conn.QueryRow(handler.ctx, statement, args...)
	} else {
		row = conn.QueryRow(handler.ctx, statement)
	}

	return row, nil
}

func (handler *PostgresHandler) Query(statement string, args ...interface{}) (basedbhandler.Rows, error) {
	var rows pgx.Rows

	conn, err := handler.pool.Acquire(handler.ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}

	if len(args) > 0 {
		rows, err = conn.Query(handler.ctx, statement, args...)
	} else {
		rows, err = conn.Query(handler.ctx, statement)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (handler *PostgresHandler) Close() {
	if handler != nil {
		handler.pool.Close()
	}
}

func NewPostgresHandler(ctx context.Context, dataSource string) (*PostgresHandler, error) {
	poolConfig, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	postgresHandler := new(PostgresHandler)
	postgresHandler.ctx = ctx
	postgresHandler.pool = pool
	//baseHandler.ErrNotFound = pgx.ErrNoRows
	return postgresHandler, nil
}
