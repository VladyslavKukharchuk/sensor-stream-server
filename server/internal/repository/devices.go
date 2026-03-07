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
	Name      string    `firestore:"name"`
	Location  string    `firestore:"location"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (d *device) toDeviceModel(id string) *model.Device {
	return &model.Device{
		ID:        id,
		MAC:       d.MAC,
		Name:      d.Name,
		Location:  d.Location,
		CreatedAt: d.CreatedAt,
	}
}

type DevicesRepository struct {
	client *firestore.Client
}

func NewDevicesRepository(client *firestore.Client) *DevicesRepository {
	return &DevicesRepository{client: client}
}

func (r *DevicesRepository) GetByMAC(ctx context.Context, mac string) (*model.Device, error) {
	iter := r.client.Collection(devicesCollection).Where("mac", "==", mac).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}

		return nil, fmt.Errorf("getting device by mac: %w", err)
	}

	var d device
	if err := doc.DataTo(&d); err != nil {
		return nil, fmt.Errorf("parsing device: %w", err)
	}

	return d.toDeviceModel(doc.Ref.ID), nil
}

func (r *DevicesRepository) Add(ctx context.Context, m *model.Device) (*model.Device, error) {
	docRef, _, err := r.client.Collection(devicesCollection).Add(ctx, device{
		MAC:       m.MAC,
		Name:      m.Name,
		Location:  m.Location,
		CreatedAt: m.CreatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("adding device: %w", err)
	}

	m.ID = docRef.ID

	return m, nil
}

func (r *DevicesRepository) List(ctx context.Context) ([]*model.Device, error) {
	iter := r.client.Collection(devicesCollection).OrderBy("created_at", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	var devices []*model.Device

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return nil, fmt.Errorf("iterating devices: %w", err)
		}

		var d device
		if err := doc.DataTo(&d); err != nil {
			return nil, fmt.Errorf("parsing device: %w", err)
		}

		devices = append(devices, d.toDeviceModel(doc.Ref.ID))
	}

	return devices, nil
}

func (r *DevicesRepository) GetByID(ctx context.Context, id string) (*model.Device, error) {
	doc, err := r.client.Collection(devicesCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting device by id: %w", err)
	}

	var d device
	if err := doc.DataTo(&d); err != nil {
		return nil, fmt.Errorf("parsing device: %w", err)
	}

	return d.toDeviceModel(doc.Ref.ID), nil
}

func (r *DevicesRepository) Update(ctx context.Context, m *model.Device) error {
	_, err := r.client.Collection(devicesCollection).Doc(m.ID).Set(ctx, device{
		MAC:       m.MAC,
		Name:      m.Name,
		Location:  m.Location,
		CreatedAt: m.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("updating device: %w", err)
	}

	return nil
}
