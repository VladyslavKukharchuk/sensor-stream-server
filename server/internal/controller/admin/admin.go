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

const (
	secondsInMinute = 60
	minutesInHour   = 60
	hoursInDay      = 24
)

type Config struct {
	FirebaseApiKey     string
	FirebaseAuthDomain string
	FirebaseProjectId  string
}

type MeasurementService interface {
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*model.Measurement, error)
	GetByDeviceID(ctx context.Context, deviceID string, since time.Time) ([]*model.Measurement, error)
}

type DevicesService interface {
	List(ctx context.Context) ([]*model.Device, error)
	GetByID(ctx context.Context, id string) (*model.Device, error)
	Update(ctx context.Context, id, name, location string) error
}

type Controller struct {
	ms     MeasurementService
	ds     DevicesService
	config Config
}

func NewController(
	ms MeasurementService,
	ds DevicesService,
	config Config,
) *Controller {
	return &Controller{
		ms:     ms,
		ds:     ds,
		config: config,
	}
}

func (c *Controller) LoginPage(f *fiber.Ctx) error {
	return f.Render("login", fiber.Map{
		"FirebaseApiKey":     c.config.FirebaseApiKey,
		"FirebaseAuthDomain": c.config.FirebaseAuthDomain,
		"FirebaseProjectId":  c.config.FirebaseProjectId,
	})
}

type DeviceDashboardItem struct {
	ID          string
	Name        string
	Location    string
	MAC         string
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

	dashboardItems := make([]DeviceDashboardItem, 0, len(devices))

	for _, device := range devices {
		item, err := c.getDeviceDashboardItem(ctx, device)
		if err != nil {
			return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		dashboardItems = append(dashboardItems, item)
	}

	return f.Render("index",
		fiber.Map{
			"Title":   "Dashboard",
			"Devices": dashboardItems,
		})
}

func (c *Controller) getDeviceDashboardItem(ctx context.Context, device *model.Device) (DeviceDashboardItem, error) {
	latest, err := c.ms.GetLatestByDeviceID(ctx, device.ID)
	if err != nil {
		return DeviceDashboardItem{}, fmt.Errorf("getting latest measurement: %w", err)
	}

	item := DeviceDashboardItem{
		ID:       device.ID,
		Name:     device.Name,
		Location: device.Location,
		MAC:      device.MAC,
	}

	if latest != nil {
		item.Temperature = latest.Temperature
		item.Humidity = latest.Humidity
		item.UpdatedAt = latest.Timestamp
		item.LastSeen = formatLastSeen(latest.Timestamp)
	} else {
		item.LastSeen = "Never"
	}

	return item, nil
}

func formatLastSeen(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < 0:
		return "Just now"
	case duration.Seconds() < secondsInMinute:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration.Minutes() < minutesInHour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration.Hours() < hoursInDay:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	default:
		return t.Format("2006-01-02 15:04")
	}
}

func (c *Controller) DevicesPage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	devices, err := c.ds.List(ctx)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	deviceItems := make([]DeviceDashboardItem, 0, len(devices))
	for _, device := range devices {
		item, err := c.getDeviceDashboardItem(ctx, device)
		if err != nil {
			return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		deviceItems = append(deviceItems, item)
	}

	return f.Render("devices", fiber.Map{
		"Title":   "Device Management",
		"Devices": deviceItems,
	})
}

type ChartData struct {
	Timestamp string  `json:"x"`
	Value     float64 `json:"y"`
}

func (c *Controller) DevicePage(f *fiber.Ctx) error {
	var (
		ctx    = context.Background()
		id     = f.Params("id")
		period = f.Query("period", "day")
	)

	device, err := c.ds.GetByID(ctx, id)
	if err != nil {
		return f.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if device == nil {
		return f.Status(http.StatusNotFound).SendString("Device not found")
	}

	var since time.Time

	switch period {
	case "month":
		since = time.Now().AddDate(0, -1, 0)
	case "week":
		since = time.Now().AddDate(0, 0, -7)
	default:
		since = time.Now().AddDate(0, 0, -1)
	}

	measurements, err := c.ms.GetByDeviceID(ctx, id, since)
	if err != nil {
		return f.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	tempData := make([]ChartData, 0, len(measurements))
	humData := make([]ChartData, 0, len(measurements))

	for _, m := range measurements {
		ts := m.Timestamp.Format(time.RFC3339)
		tempData = append(tempData, ChartData{Timestamp: ts, Value: m.Temperature})
		humData = append(humData, ChartData{Timestamp: ts, Value: m.Humidity})
	}

	return f.Render("device", fiber.Map{
		"Title":        "Device Details",
		"Device":       device,
		"TempData":     tempData,
		"HumData":      humData,
		"ActivePeriod": period,
	})
}

func (c *Controller) UpdateDevice(f *fiber.Ctx) error {
	var (
		ctx      = context.Background()
		id       = f.Params("id")
		name     = f.FormValue("name")
		location = f.FormValue("location")
	)

	if err := c.ds.Update(ctx, id, name, location); err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return f.Redirect("/admin/devices/" + id)
}
