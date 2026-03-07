package model

import "time"

type Device struct {
	ID        string
	MAC       string
	Name      string
	Location  string
	CreatedAt time.Time
}
