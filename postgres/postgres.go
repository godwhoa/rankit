package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(URI string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func IsUniqueViolation(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	return pgErr, ok && pgErr.Code == pgerrcode.UniqueViolation
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
