package userservice

import (
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

func TestGetUser(t *testing.T) {
	service.GetUser(1)
}

func TestGetUserByEmail(t *testing.T) {
	service.GetUserByEmail("test@email.com")
}

func TestCreateUser(t *testing.T) {
	var u models.User
	service.CreateUser(u, "password")
}

func TestAuthenticate(t *testing.T) {
	service.Authenticate("test@email.com", "password")
}
