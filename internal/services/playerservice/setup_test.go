package playerservice

import (
	"os"
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/repository/playerrepo"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

var service services.PlayerService

func TestMain(m *testing.M) {
	playerRepo := playerrepo.NewTestPlayerRepo()
	service = NewPlayerService(playerRepo)

	os.Exit(m.Run())
}
