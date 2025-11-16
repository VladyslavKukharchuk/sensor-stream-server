package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"sensor-stream-server/internal/model"
)

type Measurement struct {
	Temperature float64   `firestore:"temperature"`
	Humidity    float64   `firestore:"humidity"`
	Timestamp   time.Time `firestore:"timestamp"`
	CreatedAt   time.Time `firestore:"created_at"`
}

func fromMeasurementModel(m *model.Measurement) *Measurement {
	return &Measurement{
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Timestamp:   m.Timestamp,
		CreatedAt:   m.CreatedAt,
	}
}

type MeasurementRepository struct {
	client *firestore.Client
}

func NewMeasurementRepository(client *firestore.Client) *MeasurementRepository {
	return &MeasurementRepository{client: client}
}

func (r *MeasurementRepository) Add(ctx context.Context, m *model.Measurement) error {
	_, _, err := r.client.Collection("measurements").Add(ctx, fromMeasurementModel(m))
	if err != nil {
		return fmt.Errorf("failed to add measurement in firestore: %w", err)
	}

	return nil
}
