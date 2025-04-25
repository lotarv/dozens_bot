package storage

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ErrInvalidTxType = errors.New("invalid transaction type")
)

type TxKey struct{}
type Storage struct {
	db *sqlx.DB
}

func (s *Storage) DB() *sqlx.DB {
	return s.db
}

func New() (*Storage, error) {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	slog.Info("Connection string: %s", connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
