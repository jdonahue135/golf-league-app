package repository

import (
	"context"
	"database/sql"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

type UserRepo interface {
	CreateUser(u models.User, password string) (int, error)
	Authenticate(email, password string) (int, int, error)
	AllUsers() bool
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(u models.User) error
	BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error)
	CreateInactiveUserTransaction(u models.User, ctx context.Context, tx *sql.Tx) (int, error)
}
