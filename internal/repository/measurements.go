package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"sensor-stream-server/internal/model"
)

type measurement struct {
	DeviceID    string    `firestore:"device_id"`
	Temperature float64   `firestore:"temperature"`
	Humidity    float64   `firestore:"humidity"`
	Timestamp   time.Time `firestore:"timestamp"`
	CreatedAt   time.Time `firestore:"created_at"`
}

func (m *measurement) toMeasurementModel() *model.Measurement {
	return &model.Measurement{
		DeviceID:    m.DeviceID,
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Timestamp:   m.Timestamp,
		CreatedAt:   m.CreatedAt,
	}
}

func fromMeasurementModel(m *model.Measurement) *measurement {
	return &measurement{
		DeviceID:    m.DeviceID,
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

func (r *MeasurementRepository) List(ctx context.Context) ([]*model.Measurement, error) {
	iter := r.client.Collection("measurements").OrderBy("created_at", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	var measurements []*model.Measurement

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return nil, fmt.Errorf("failed to iterate measurements: %w", err)
		}

		var m measurement
		if err := doc.DataTo(&m); err != nil {
			return nil, fmt.Errorf("failed to parse measurement document: %w", err)
		}

		measurements = append(measurements, m.toMeasurementModel())
	}

	return measurements, nil
}
