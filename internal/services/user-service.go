package services

import "github.com/jdonahue135/golf-league-app/internal/models"

type UserService interface {
	GetUser(userID int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User, password string) (int, error)
	Authenticate(email, password string) (int, int, error)
}
