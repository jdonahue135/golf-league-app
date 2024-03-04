package repository

import (
	"github.com/jdonahue135/golf-league-app/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(u models.User) error
	GetLeagueByName(name string) (models.League, error)
	GetLeagueByID(id int) (models.League, error)
	CreateLeague(league models.League) (int, error)
	Authenticate(email, password string) (int, int, error)
	CreateUser(u models.User, password string) (int, error)
}
