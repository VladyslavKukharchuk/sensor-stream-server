//nolint:wrapcheck
package admin

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type MeasurementService interface {
	List(ctx context.Context) ([]*model.Measurement, error)
}

type DevicesService interface {
	List(ctx context.Context) ([]*model.Device, error)
	GetByID(ctx context.Context, id string) (*model.Device, error)
}

type Controller struct {
	ms MeasurementService
	ds DevicesService
}

func NewController(
	ms MeasurementService,
	ds DevicesService,
) *Controller {
	return &Controller{
		ms: ms,
		ds: ds,
	}
}

func (c *Controller) IndexPage(f *fiber.Ctx) error {
	return f.Render("index",
		fiber.Map{
			"Title": "Main",
		})
}

func (c *Controller) MeasurementsPage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	measurements, err := c.ms.List(ctx)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return f.Render("measurements", fiber.Map{
		"Title":        "Measurements",
		"Measurements": measurements,
	})
}

func (c *Controller) DevicesPage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	devices, err := c.ds.List(ctx)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return f.Render("devices", fiber.Map{
		"Title":   "Devices",
		"Devices": devices,
	})
}

func (c *Controller) DevicePage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
		id  = f.Params("id")
	)

	device, err := c.ds.GetByID(ctx, id)
	if err != nil {
		return f.Status(500).SendString(err.Error())
	}

	if device == nil {
		return f.Status(404).SendString("Device not found")
	}

	return f.Render("device", fiber.Map{
		"Title":  "Device " + id,
		"Device": device,
	})
}
