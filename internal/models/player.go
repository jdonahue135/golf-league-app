package models

import (
	"time"
)

type Player struct {
	ID             int
	LeagueID       int
	UserID         int
	Handicap       int
	IsCommissioner bool
	IsActive       bool
	User           User
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
