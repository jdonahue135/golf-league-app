package models

import (
	"time"
)

// LeagueAdmin is the league admin model
type LeagueAdmin struct {
	ID        int
	LeagueID  int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
