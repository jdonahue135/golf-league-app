package playerservice

import (
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

func TestGetPlayersInLeague(t *testing.T) {
	service.GetPlayersInLeague(1)
}

func TestGetPlayerInLeague(t *testing.T) {
	service.GetPlayerInLeague(1, 1)
}
func TestActivatePlayer(t *testing.T) {
	var p models.Player
	service.ActivatePlayer(p)
}
