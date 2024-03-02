package models

import (
	"time"
)

const (
	AccessLevelPlayer = iota + 1
	AccessLevelAdmin
	AccessLevelSuperAdmin
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
