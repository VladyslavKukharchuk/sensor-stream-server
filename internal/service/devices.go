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
	List(ctx context.Context) ([]*model.Device, error)
	GetByID(ctx context.Context, id string) (*model.Device, error)
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

func (s *DevicesService) List(ctx context.Context) ([]*model.Device, error) {
	devices, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting devices: %w", err)
	}

	return devices, nil
}

func (s *DevicesService) GetByID(ctx context.Context, id string) (*model.Device, error) {
	device, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting device: %w", err)
	}

	return device, nil
}
