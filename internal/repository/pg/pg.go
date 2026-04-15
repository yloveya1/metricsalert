package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

const (
	UpdateMetricQuery = `
        INSERT INTO metrics (name, type, delta, value)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (name) DO UPDATE SET
            type  = EXCLUDED.type,
            delta = metrics.delta + EXCLUDED.delta,
            value = EXCLUDED.value,
            updated_at = NOW();`
)

type Database struct {
	pg *pgxpool.Pool
}

func NewDatabase(ctx context.Context, connStr string) (repository.IStorage, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to open sql.DB for migrations: %w", err)
	}
	defer db.Close()

	migrationsDir := "file://migrations"
	if err = runMigrations(db, migrationsDir); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Database{pg: pool}, nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsDir,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func (db *Database) Close() {
	db.pg.Close()
}

func (db *Database) UpdateCounterMetric(ctx context.Context, metric *models.Metrics) error {
	return db.updateMetric(ctx, metric)

}

func (db *Database) UpdateGaugeMetric(ctx context.Context, metric *models.Metrics) error {
	return db.updateMetric(ctx, metric)
}

func (db *Database) GetMetricList(ctx context.Context) ([]*models.Metrics, error) {
	rows, err := db.pg.Query(ctx, `SELECT name, type, delta, value FROM metrics`)
	if err != nil {
		return nil, fmt.Errorf("failed to query metric list: %w", err)
	}
	defer rows.Close()

	var metrics []*models.Metrics
	for rows.Next() {
		m := &models.Metrics{}
		err = rows.Scan(
			&m.ID,
			&m.MType,
			&m.Delta,
			&m.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metric row: %w", err)
		}

		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan metric rows: %w", err)
	}

	return metrics, nil
}

func (db *Database) GetMetricByID(ctx context.Context, metric *models.Metrics) (models.Metrics, error) {
	m := models.Metrics{}
	err := db.pg.QueryRow(ctx, `SELECT name, type, delta, value FROM metrics where name = $1`, metric.ID).Scan(
		&m.ID,
		&m.MType,
		&m.Delta,
		&m.Value,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metrics{}, metrics.ErrMetricNotFound
		}

		return models.Metrics{}, fmt.Errorf("failed to send query: %w", err)
	}

	return m, nil
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.pg.Ping(ctx)
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}

func (db *Database) updateMetric(ctx context.Context, metric *models.Metrics) error {
	_, err := db.pg.Exec(ctx, UpdateMetricQuery,
		metric.ID, metric.MType, metric.Delta, metric.Value)

	if err != nil {
		return fmt.Errorf("failed to update %s metric: %w", metric.MType, err)
	}

	return nil
}

func (db *Database) UpdateMetricList(ctx context.Context, metrics []*models.Metrics) error {
	tx, err := db.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx, err: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, metric := range metrics {
		_, err = tx.Exec(ctx, UpdateMetricQuery, metric.ID, metric.MType, metric.Delta, metric.Value)
		if err != nil {
			return fmt.Errorf("failed to update metric: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit tx, err: %w", err)
	}

	return nil
}
