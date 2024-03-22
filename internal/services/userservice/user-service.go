package userservice

import (
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type userService struct {
	UserRepo repository.UserRepo
}

func NewUserService(r repository.UserRepo) services.UserService {
	return &userService{UserRepo: r}
}

func (m *userService) GetUser(userID int) (models.User, error) {
	return m.UserRepo.GetUserByID(userID)
}

func (m *userService) GetUserByEmail(email string) (models.User, error) {
	return m.UserRepo.GetUserByEmail(email)
}

func (m *userService) CreateUser(user models.User, password string) (int, error) {
	return m.UserRepo.CreateUser(user, password)
}

func (m *userService) Authenticate(email, password string) (int, int, error) {
	return m.UserRepo.Authenticate(email, password)
}
