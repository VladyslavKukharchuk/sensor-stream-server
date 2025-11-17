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

const devicesCollection = "devices"

type device struct {
	MAC       string    `firestore:"mac"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (m *device) toDeviceModel(id string) *model.Device {
	return &model.Device{
		ID:        id,
		MAC:       m.MAC,
		CreatedAt: m.CreatedAt,
	}
}

type DevicesRepository struct {
	client *firestore.Client
}

func NewDevicesRepository(client *firestore.Client) *DevicesRepository {
	return &DevicesRepository{client: client}
}

func (r *DevicesRepository) GetByMAC(ctx context.Context, mac string) (*model.Device, error) {
	iter := r.client.Collection(devicesCollection).
		Where("mac", "==", mac).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByMAC firestore query: %w", err)
	}

	var d device
	if err := doc.DataTo(&d); err != nil {
		return nil, fmt.Errorf("parsing device data: %w", err)
	}

	return d.toDeviceModel(doc.Ref.ID), nil
}

func (r *DevicesRepository) Add(ctx context.Context, m *model.Device) (*model.Device, error) {
	docRef := r.client.Collection(devicesCollection).NewDoc()

	data := device{
		MAC:       m.MAC,
		CreatedAt: time.Now().UTC(),
	}

	_, err := docRef.Set(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("adding device to firestore: %w", err)
	}

	m.ID = docRef.ID
	m.CreatedAt = data.CreatedAt

	return m, nil
}
