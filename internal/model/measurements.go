package model

import "time"

type Measurement struct {
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}
