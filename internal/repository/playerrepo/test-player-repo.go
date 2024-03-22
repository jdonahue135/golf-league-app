package playerrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type testPlayerRepo struct{}

func NewTestPlayerRepo() repository.PlayerRepo {
	return &testPlayerRepo{}
}

func (m *testPlayerRepo) CreatePlayer(player models.Player) error {
	if player.LeagueID == 1 {
		return errors.New("db error")
	}
	return nil
}

func (m *testPlayerRepo) UpdatePlayer(p models.Player) error {
	if p.UserID == 2 {
		return errors.New("some error")
	}
	return nil
}

func (m *testPlayerRepo) GetPlayersByLeagueID(leagueID int) ([]models.Player, error) {
	if leagueID == 2 {
		return nil, errors.New("some error")
	}
	var p []models.Player
	return p, nil
}

func (m *testPlayerRepo) GetPlayerByUserAndLeagueID(userID, leagueID int) (models.Player, error) {
	var p models.Player
	if userID == 0 {
		return p, errors.New("some error")
	}
	if userID == 1 {
		p.IsCommissioner = true
	} else {
		p.IsCommissioner = false
	}
	if userID == 2 || userID == 4 {
		p.UserID = userID
		p.IsActive = false
	}
	if userID == 3 {
		p.IsActive = true
	}
	return p, nil
}

func (m *testPlayerRepo) CreatePlayerTransaction(player models.Player, ctx context.Context, tx *sql.Tx) error {
	if player.Handicap == 100 {
		return errors.New("player error")
	}
	if player.UserID == 2 {
		return errors.New("player error")
	}
	return nil
}

func (m *testPlayerRepo) CommitTransaction(tx *sql.Tx) error {
	return nil
}
