package models

import (
	"time"
)

// League is the league model
type League struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
