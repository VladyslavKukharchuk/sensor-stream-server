//nolint:wrapcheck
package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type MeasurementService interface {
	List(ctx context.Context) ([]*model.Measurement, error)
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*model.Measurement, error)
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

type DeviceDashboardItem struct {
	ID          string
	Temperature float64
	Humidity    float64
	LastSeen    string
	UpdatedAt   time.Time
}

func (c *Controller) IndexPage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	devices, err := c.ds.List(ctx)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var dashboardItems []DeviceDashboardItem

	for _, device := range devices {
		latest, err := c.ms.GetLatestByDeviceID(ctx, device.ID)
		if err != nil {
			return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		item := DeviceDashboardItem{
			ID: device.ID,
		}

		if latest != nil {
			item.Temperature = latest.Temperature
			item.Humidity = latest.Humidity
			item.UpdatedAt = latest.Timestamp
			item.LastSeen = formatLastSeen(latest.Timestamp)
		} else {
			item.LastSeen = "Never"
		}

		dashboardItems = append(dashboardItems, item)
	}

	return f.Render("index",
		fiber.Map{
			"Title":   "Dashboard",
			"Devices": dashboardItems,
		})
}

func formatLastSeen(t time.Time) string {
	duration := time.Since(t)
	if duration.Seconds() < 60 {
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	}
	if duration.Minutes() < 60 {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	}
	if duration.Hours() < 24 {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	}
	return t.Format("2006-01-02 15:04")
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
		return f.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if device == nil {
		return f.Status(http.StatusNotFound).SendString("Device not found")
	}

	return f.Render("device", fiber.Map{
		"Title":  "Device " + id,
		"Device": device,
	})
}
