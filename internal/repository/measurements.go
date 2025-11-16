package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"sensor-stream-server/internal/model"
)

type Measurement struct {
	Temperature float64   `firestore:"temperature"`
	Humidity    float64   `firestore:"humidity"`
	Timestamp   time.Time `firestore:"timestamp"`
}

func fromMeasurementModel(m *model.Measurement) *Measurement {
	return &Measurement{
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Timestamp:   m.Timestamp,
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

	return err
}
