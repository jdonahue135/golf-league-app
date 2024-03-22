package leaguerepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type testLeagueRepo struct{}

func NewTestLeagueRepo() repository.LeagueRepo {
	return &testLeagueRepo{}
}

func (m *testLeagueRepo) BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error) {
	return context.Background(), func() {}, &sql.Tx{}, nil
}

func (m *testLeagueRepo) CreateLeagueTransaction(league models.League, ctx context.Context, tx *sql.Tx) (int, error) {
	if league.Name == "League Error" {
		return 0, errors.New("league creation failed")
	}
	league.ID = 1
	return 1, nil
}

func (m *testLeagueRepo) GetLeagueByName(name string) (models.League, error) {
	var l models.League
	if name == "league1" || name == "league2" {
		return l, errors.New("some error")
	}
	return l, nil
}

// GetLeagueByID returns a league by ID
func (m *testLeagueRepo) GetLeagueByID(id int) (models.League, error) {
	var l models.League

	if id == 3 {
		return l, errors.New("some error")
	}

	l.ID = id
	return l, nil
}

func (m *testLeagueRepo) CreateLeague(league models.League, commissioner models.Player) (int, error) {
	if league.Name == "league1" {
		return 0, errors.New("some error")
	}
	return 1, nil
}
