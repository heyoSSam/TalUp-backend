package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

func GetDBConnection() (*pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
