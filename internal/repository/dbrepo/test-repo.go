package dbrepo

import (
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	if id == 0 {
		return u, errors.New("some error")
	}
	return u, nil
}

func (m *testDBRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	if email == "me@here.ca" {
		return u, errors.New("some error")
	}
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) CreateUser(u models.User, password string) (int, error) {
	if password == "error" {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) GetLeagueByName(name string) (models.League, error) {
	var l models.League
	if name == "league1" || name == "league2" {
		return l, errors.New("some error")
	}
	return l, nil
}

// GetLeagueByID returns a league by ID
func (m *testDBRepo) GetLeagueByID(id int) (models.League, error) {
	var l models.League

	if id == 3 {
		return l, errors.New("some error")
	}

	l.ID = id
	return l, nil
}

func (m *testDBRepo) CreateLeague(league models.League, commissioner models.Player) (int, error) {
	if league.Name == "league1" {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) GetPlayersByLeagueID(leagueID int) ([]models.Player, error) {
	if leagueID == 2 {
		return nil, errors.New("some error")
	}
	var p []models.Player
	return p, nil
}

func (m *testDBRepo) Authenticate(email, password string) (int, int, error) {
	if email == "jack@nimble.com" {
		return 0, 0, errors.New("some error")
	}
	return 1, 1, nil
}
