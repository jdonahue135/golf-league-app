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

func TestGetPlayer(t *testing.T) {
	service.GetPlayer(1)
}

func TestActivatePlayer(t *testing.T) {
	var p models.Player
	service.ActivatePlayer(p)
}

func TestRemovePlayer(t *testing.T) {
	var p models.Player
	err := service.RemovePlayer(p)
	if err != nil {
		t.Error("failed success: expected no error but got one")
	}
	p.IsCommissioner = true
	err = service.RemovePlayer(p)
	if err == nil {
		t.Error("failed error: expected error but got none")
	}
}
