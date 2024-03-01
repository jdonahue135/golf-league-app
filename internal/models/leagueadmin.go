package models

import (
	"time"
)

// League is the league model
type LeagueAdmin struct {
	ID        int
	LeagueID  int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
