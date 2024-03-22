package leagueservice

import (
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type testLeagueService struct {
	LeagueRepo repository.LeagueRepo
	PlayerRepo repository.PlayerRepo
	UserRepo   repository.UserRepo
}

func NewTestLeagueService(l repository.LeagueRepo, p repository.PlayerRepo, u repository.UserRepo) services.LeagueService {
	return &testLeagueService{LeagueRepo: l, PlayerRepo: p, UserRepo: u}
}

func (m *testLeagueService) GetLeague(ID int) (models.League, error) {
	var l models.League
	if ID == 3 {
		return l, errors.New("league doesn't exist")
	}
	l.ID = ID
	return l, nil
}

func (m *testLeagueService) GetLeagueByName(name string) (models.League, error) {
	var l models.League
	if name == "league0" {
		return l, nil
	}
	return l, errors.New("league name not found in DB")
}

func (m *testLeagueService) CreateLeagueWithCommissioner(league models.League, commissioner models.Player) (int, error) {
	if league.Name == "league1" {
		return 0, errors.New("error inserting league in DB")
	}
	return 1, nil
}

func (m *testLeagueService) AddExistingUserToLeague(userID, leagueID int) error {
	if leagueID == 6 {
		return errors.New("user already active in league")
	}
	return nil
}

func (m *testLeagueService) AddNewUserToLeague(user models.User, leagueID int) error {
	if leagueID == 2 {
		return errors.New("error adding user to DB")
	}
	return nil
}
