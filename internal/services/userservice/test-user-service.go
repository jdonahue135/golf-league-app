package userservice

import (
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type testUserService struct {
	UserRepo repository.UserRepo
}

func NewTestUserService(r repository.UserRepo) services.UserService {
	return &testUserService{UserRepo: r}
}

func (m *testUserService) GetUser(userID int) (models.User, error) {
	var u models.User
	if userID == 0 {
		return u, errors.New("user not found")
	}
	u.ID = userID

	return u, nil
}

func (m *testUserService) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	if email == "me@here.ca" {
		return u, errors.New("user not found")
	}

	return u, nil
}

func (m *testUserService) CreateUser(user models.User, password string) (int, error) {
	if password == "error" {
		return 0, errors.New("sign up error")
	}
	return 1, nil
}

func (m *testUserService) Authenticate(email, password string) (int, int, error) {
	if email == "jack@nimble.com" {
		return 0, 0, errors.New("Invalid credentials")
	}
	return 1, 1, nil
}
