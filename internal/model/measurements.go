package model

import "time"

type Measurement struct {
	DeviceID    string
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
	CreatedAt   time.Time
}
