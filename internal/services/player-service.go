package services

import "github.com/jdonahue135/golf-league-app/internal/models"

type PlayerService interface {
	GetPlayer(ID int) (models.Player, error)
	GetPlayersInLeague(leagueID int) ([]models.Player, error)
	GetPlayerInLeague(userID, leagueID int) (models.Player, error)
	ActivatePlayer(player models.Player) error
	RemovePlayer(player models.Player) error
}
