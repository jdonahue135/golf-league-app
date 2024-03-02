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
	return u, nil
}

func (m *testDBRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) CreateUser(u models.User, password string) (int, error) {
	return 1, nil
}

func (m *testDBRepo) GetLeagueByName(name string) (models.League, error) {
	var l models.League
	return l, nil
}

// GetLeagueByID returns a league by ID
func (m *testDBRepo) GetLeagueByID(id int) (models.League, error) {
	var l models.League
	return l, nil
}

func (m *testDBRepo) CreateLeague(league models.League) (int, error) {
	return 1, nil
}

func (m *testDBRepo) Authenticate(email, password string) (int, error) {
	if email == "jack@nimble.com" {
		return 0, errors.New("some error")
	}
	return 1, nil
}
