package services

import "github.com/jdonahue135/golf-league-app/internal/models"

type PlayerService interface {
	GetPlayersInLeague(leagueID int) ([]models.Player, error)
	GetPlayerInLeague(userID, leagueID int) (models.Player, error)
	ActivatePlayer(player models.Player) error
}
