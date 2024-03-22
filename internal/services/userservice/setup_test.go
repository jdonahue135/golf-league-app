package userservice

import (
	"os"
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/repository/userrepo"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

var service services.UserService

func TestMain(m *testing.M) {
	userRepo := userrepo.NewTestUserRepo()
	service = NewUserService(userRepo)

	os.Exit(m.Run())
}
