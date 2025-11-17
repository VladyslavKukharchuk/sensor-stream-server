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
	ID        string    `firestore:"id"`
	MAC       string    `firestore:"mac"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (m *device) toDeviceModel() *model.Device {
	return &model.Device{
		ID:        m.ID,
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

	d.ID = doc.Ref.ID

	return d.toDeviceModel(), nil
}

func (r *DevicesRepository) Add(ctx context.Context, device *model.Device) (*model.Device, error) {
	docRef := r.client.Collection(devicesCollection).NewDoc()
	device.ID = docRef.ID

	_, err := docRef.Set(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("adding device to firestore: %w", err)
	}

	return device, nil
}
