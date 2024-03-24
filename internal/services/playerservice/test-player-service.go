package playerservice

import (
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type testPlayerService struct {
	PlayerRepo repository.PlayerRepo
}

func NewTestPlayerService(r repository.PlayerRepo) services.PlayerService {
	return &testPlayerService{PlayerRepo: r}
}

func (m *testPlayerService) GetPlayer(ID int) (models.Player, error) {
	var p models.Player
	if ID == 9 {
		return p, errors.New("Player not found")
	}
	if ID == 8 {
		p.IsActive = false
	}
	p.ID = ID
	p.IsActive = true
	return p, nil
}

func (m *testPlayerService) GetPlayersInLeague(leagueID int) ([]models.Player, error) {
	var p []models.Player
	if leagueID == 2 {
		return p, errors.New("player error")
	}
	return p, nil
}

func (m *testPlayerService) GetPlayerInLeague(userID, leagueID int) (models.Player, error) {
	var p models.Player
	if userID == 4 && leagueID == 4 {
		return p, errors.New("user not in league")
	}
	if userID == 3 {
		p.IsCommissioner = false
	} else {
		p.IsCommissioner = true
	}
	return p, nil
}

func (m *testPlayerService) ActivatePlayer(player models.Player) error {
	return nil
}

func (m *testPlayerService) RemovePlayer(player models.Player) error {
	if player.ID == 10 {
		return errors.New("service error")
	}
	return nil
}
