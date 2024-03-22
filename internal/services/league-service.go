package services

import "github.com/jdonahue135/golf-league-app/internal/models"

type LeagueService interface {
	GetLeague(ID int) (models.League, error)
	GetLeagueByName(name string) (models.League, error)
	CreateLeagueWithCommissioner(league models.League, commissioner models.Player) (int, error)
	AddExistingUserToLeague(userID, leagueID int) error
	AddNewUserToLeague(user models.User, leagueID int) error
}
