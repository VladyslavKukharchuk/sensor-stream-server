package service

import (
	"context"
	"fmt"

	"sensor-stream-server/internal/model"
)

type Repository interface {
	Add(ctx context.Context, m model.Measurement) error
}

type MeasurementService struct {
	repository Repository
}

func NewMeasurementService(repository Repository) *MeasurementService {
	return &MeasurementService{repository: repository}
}

func (s *MeasurementService) Add(ctx context.Context, m model.Measurement) error {
	if err := s.repository.Add(ctx, m); err != nil {
		return fmt.Errorf("adding measurement: %w", err)
	}

	return nil
}
