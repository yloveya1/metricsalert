package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
)

type Database struct {
	pg *pgxpool.Pool
}

func (db *Database) UpdateCounterMetric(metric *models.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (db *Database) UpdateGaugeMetric(metric *models.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (db *Database) GetMetricList() ([]*models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) GetMetricByID(metric *models.Metrics) (models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.pg.Ping(ctx)
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}

func NewDatabase(ctx context.Context, connStr string) (repository.IStorage, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	return &Database{pg: pool}, nil
}

func (db *Database) Close() {
	db.pg.Close()
}
