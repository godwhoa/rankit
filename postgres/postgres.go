package postgres

import (
	"context"
	"errors"
	"fmt"

	"rankit/postgres/sqlgen"

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

type Querier interface {
	sqlgen.Querier
	WithTx(pgx.Tx) sqlgen.Querier
}

func NewQuerier(pool *pgxpool.Pool) *sqlgen.Queries {
	return sqlgen.New(pool)
}

func IsUniqueViolation(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	return pgErr, ok && pgErr.Code == pgerrcode.UniqueViolation
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsForeignKeyViolation(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	return pgErr, ok && pgErr.Code == pgerrcode.ForeignKeyViolation
}
