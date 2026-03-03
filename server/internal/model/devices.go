package model

import "time"

type Device struct {
	ID        string
	MAC       string
	CreatedAt time.Time
}
