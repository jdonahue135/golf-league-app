package models

import "time"

type Player struct {
	ID             int
	LeagueID       int
	UserID         int
	Handicap       int
	IsCommissioner bool
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
