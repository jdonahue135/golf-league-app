package playerservice

import (
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type playerService struct {
	PlayerRepo repository.PlayerRepo
}

func NewPlayerService(r repository.PlayerRepo) services.PlayerService {
	return &playerService{PlayerRepo: r}
}

func (m *playerService) GetPlayersInLeague(leagueID int) ([]models.Player, error) {
	return m.PlayerRepo.GetPlayersByLeagueID(leagueID)
}

func (m *playerService) GetPlayerInLeague(userID, leagueID int) (models.Player, error) {
	return m.PlayerRepo.GetPlayerByUserAndLeagueID(userID, leagueID)
}

func (m *playerService) ActivatePlayer(player models.Player) error {
	player.IsActive = true
	return m.PlayerRepo.UpdatePlayer(player)
}
