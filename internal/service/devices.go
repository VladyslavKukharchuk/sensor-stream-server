package service

import (
	"context"
	"fmt"
	"time"

	"sensor-stream-server/internal/model"
)

type DevicesRepository interface {
	GetByMAC(context.Context, string) (*model.Device, error)
	Add(ctx context.Context, m *model.Device) (*model.Device, error)
}

type DevicesService struct {
	repository DevicesRepository
}

func NewDevicesService(repository DevicesRepository) *DevicesService {
	return &DevicesService{repository: repository}
}

func (s *DevicesService) Add(ctx context.Context, mac string) (*model.Device, error) {
	existingDevice, err := s.repository.GetByMAC(ctx, mac)
	if err == nil && existingDevice != nil {
		return existingDevice, nil
	}

	device := &model.Device{
		MAC:       mac,
		CreatedAt: time.Now().UTC(),
	}
	newDevice, err := s.repository.Add(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("adding device: %w", err)
	}

	return newDevice, nil
}
