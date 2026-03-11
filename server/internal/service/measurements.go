package service

import (
	"context"
	"fmt"
	"sort"
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

func (s *MeasurementService) GetAggregatedByDeviceID(ctx context.Context, deviceID string, since time.Time, interval time.Duration) ([]*model.Measurement, error) {
	measurements, err := s.repository.GetByDeviceID(ctx, deviceID, since)
	if err != nil {
		return nil, fmt.Errorf("getting measurements for device %s: %w", deviceID, err)
	}

	if len(measurements) == 0 {
		return nil, nil
	}

	if interval <= 0 {
		return measurements, nil
	}

	type aggData struct {
		tempSum float64
		humSum  float64
		count   int
	}

	buckets := make(map[time.Time]*aggData)

	for _, m := range measurements {
		// Group by interval window
		bucketTime := m.Timestamp.Truncate(interval)
		if _, ok := buckets[bucketTime]; !ok {
			buckets[bucketTime] = &aggData{}
		}
		buckets[bucketTime].tempSum += m.Temperature
		buckets[bucketTime].humSum += m.Humidity
		buckets[bucketTime].count++
	}

	result := make([]*model.Measurement, 0, len(buckets))
	for t, data := range buckets {
		result = append(result, &model.Measurement{
			DeviceID:    deviceID,
			Timestamp:   t,
			Temperature: float64(int(data.tempSum/float64(data.count)*10)) / 10,
			Humidity:    float64(int(data.humSum/float64(data.count)*10)) / 10,
		})
	}

	// Sorting by timestamp is important for charts
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result, nil
}
