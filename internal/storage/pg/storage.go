package pg

import (
	"context"
	"fmt"
	"github.com/MikhailSolovev/URLAsker/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const noRows = "no rows in result set"

type Storage struct {
	pool *pgxpool.Pool
}

// NewPool TODO: refactor for connection string
func NewPool(ctx context.Context, user, pass, host, port, db string, timeout time.Duration) (*Storage, *pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, pass, host, port, db))
	if err != nil {
		return nil, nil, err
	}

	return &Storage{pool: pool}, pool, nil
}

func (s *Storage) ListLatestResult(ctx context.Context) (results models.ResultsPostgresDTO, err error) {
	q := `SELECT dtm, url, available FROM availability
		  WHERE dtm = (SELECT dtm FROM availability 
		               ORDER BY dtm DESC
		               LIMIT 1)`

	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		if err.Error() == noRows {
			return results, nil
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return results, pgErr
		}
		return
	}

	var result models.ResultPostgresDTO

	for rows.Next() {
		if err = rows.Scan(&result.Date, &result.URL, &result.Available); err != nil {
			return results, err
		}

		results.Results = append(results.Results, result)
	}

	if err = rows.Err(); err != nil {
		return results, err
	}

	return results, nil
}

// ListResults TODO: make pagination (offset and limit)
func (s *Storage) ListResults(ctx context.Context, dateFrom, dateTo time.Time) (results models.ResultsPostgresDTO, err error) {
	q := `SELECT dtm, url, available FROM availability
		  WHERE dtm BETWEEN $1 AND $2
		  ORDER BY dtm DESC`

	rows, err := s.pool.Query(ctx, q, dateFrom.UTC(), dateTo.UTC())
	if err != nil {
		if err.Error() == noRows {
			return results, nil
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return results, pgErr
		}
		return
	}

	var result models.ResultPostgresDTO

	for rows.Next() {
		if err = rows.Scan(&result.Date, &result.URL, &result.Available); err != nil {
			return results, err
		}

		results.Results = append(results.Results, result)
	}

	if err = rows.Err(); err != nil {
		return results, err
	}

	return results, nil
}

func (s *Storage) RecordResult(ctx context.Context, result models.ResultPostgresDTO) (err error) {
	q := `INSERT INTO availability (dtm, url, available)
		  VALUES ($1, $2, $3)`

	if err = s.pool.QueryRow(ctx, q, result.Date, result.URL, result.Available).Scan(); err != nil {
		if err.Error() == noRows {
			return nil
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return pgErr
		}
		return
	}

	return
}
