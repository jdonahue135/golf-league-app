package repository

import (
	"context"
	"database/sql"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

type LeagueRepo interface {
	GetLeagueByName(name string) (models.League, error)
	GetLeagueByID(id int) (models.League, error)
	GetLeaguesByUserID(userID int) ([]models.League, error)
	CreateLeague(league models.League, commissioner models.Player) (int, error)
	CreateLeagueTransaction(league models.League, ctx context.Context, tx *sql.Tx) (int, error)
}
