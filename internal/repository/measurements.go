package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sensor-stream-server/internal/model"
)

type Measurement struct {
	ID          int
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

type MeasurementRepository struct {
	db *pgxpool.Pool
}

func NewMeasurementRepository(db *pgxpool.Pool) *MeasurementRepository {
	return &MeasurementRepository{db: db}
}

func (r *MeasurementRepository) Add(ctx context.Context, m model.Measurement) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO measurements (temperature, humidity, timestamp) VALUES ($1, $2, $3)",
		m.Temperature, m.Humidity, m.Timestamp,
	)

	return err
}
