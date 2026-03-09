package service

import (
	"context"
	"fmt"
	"time"

	"sensor-stream-server/internal/model"
)

type Repository interface {
	Add(ctx context.Context, m *model.Measurement) error
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*model.Measurement, error)
	GetByDeviceID(ctx context.Context, deviceID string, since time.Time) ([]*model.Measurement, error)
}

type MeasurementService struct {
	repository Repository
}

func NewMeasurementService(repository Repository) *MeasurementService {
	return &MeasurementService{repository: repository}
}

func (s *MeasurementService) Add(ctx context.Context, m *model.Measurement) error {
	m.CreatedAt = time.Now().UTC()
	if err := s.repository.Add(ctx, m); err != nil {
		return fmt.Errorf("adding measurement: %w", err)
	}

	return nil
}

func (s *MeasurementService) GetLatestByDeviceID(ctx context.Context, deviceID string) (*model.Measurement, error) {
	measurement, err := s.repository.GetLatestByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("getting latest measurement for device %s: %w", deviceID, err)
	}

	return measurement, nil
}

func (s *MeasurementService) GetByDeviceID(ctx context.Context, deviceID string, since time.Time) ([]*model.Measurement, error) {
	measurements, err := s.repository.GetByDeviceID(ctx, deviceID, since)
	if err != nil {
		return nil, fmt.Errorf("getting measurements for device %s: %w", deviceID, err)
	}

	return measurements, nil
}
