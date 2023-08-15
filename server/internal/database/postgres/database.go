package postgres

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// need migration.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dsn string, maxAttempts int) (db *Repository, err error) {
	var pool *pgxpool.Pool
	const pause = 5 * time.Second
	err = doWithTries(func() error {
		ctxCancel, cancel := context.WithTimeout(ctx, pause)
		defer cancel()

		pool, err = pgxpool.Connect(ctxCancel, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, pause)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	ok, err := startMigration(dsn)
	if err != nil && !ok {
		return nil, fmt.Errorf("failed migrate database: %w", err)
	}

	return &Repository{db: pool}, err
}

func startMigration(dsn string) (bool, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}
	return true, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}
