package repository

import (
	"context"
	"database/sql"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

type PlayerRepo interface {
	CreatePlayer(player models.Player) error
	UpdatePlayer(p models.Player) error
	GetPlayerByID(ID int) (models.Player, error)
	GetPlayersByLeagueID(leagueID int) ([]models.Player, error)
	GetPlayerByUserAndLeagueID(userID, leagueID int) (models.Player, error)
	CreatePlayerTransaction(player models.Player, ctx context.Context, tx *sql.Tx) error
}
