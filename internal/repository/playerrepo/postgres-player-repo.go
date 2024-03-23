package playerrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type postgresPlayerRepo struct {
	DB *sql.DB
}

func NewPostgresPlayerRepo(conn *sql.DB) repository.PlayerRepo {
	return &postgresPlayerRepo{
		DB: conn,
	}
}

// UpdatePlayer updates a player in the db
func (m *postgresPlayerRepo) UpdatePlayer(p models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update players set handicap = $1, is_commissioner = $2, isActive = $3, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		p.Handicap,
		p.IsCommissioner,
		p.IsActive,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresPlayerRepo) CreatePlayer(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into players 
		(league_id, user_id, handicap, is_commissioner, is_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		player.LeagueID,
		player.UserID,
		player.Handicap,
		player.IsCommissioner,
		player.IsActive,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresPlayerRepo) GetPlayerByUserAndLeagueID(userID, leagueID int) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		id,
		league_id,
		user_id,
		is_commissioner,
		is_active,
		created_at,
		updated_at
	from players 
	where league_id=$1 and user_id = $2`

	row := m.DB.QueryRowContext(ctx, query, leagueID, userID)

	var p models.Player

	err := row.Scan(
		&p.ID,
		&p.LeagueID,
		&p.UserID,
		&p.IsCommissioner,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (m *postgresPlayerRepo) GetPlayersByLeagueID(leagueID int) ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		p.id,
		p.league_id,
		p.user_id,
		p.is_commissioner,
		p.is_active,
		p.created_at,
		p.updated_at,
		u.id,
		u.first_name,
		u.last_name
	from players p join users u on p.user_id = u.id 
	where league_id=$1`

	var players []models.Player

	rows, err := m.DB.QueryContext(ctx, query, leagueID)
	if err != nil {
		return players, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Player

		err := rows.Scan(
			&p.ID,
			&p.LeagueID,
			&p.UserID,
			&p.IsCommissioner,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.User.ID,
			&p.User.FirstName,
			&p.User.LastName,
		)
		if err != nil {
			return players, err
		}

		players = append(players, p)
	}

	if err = rows.Err(); err != nil {
		return players, err
	}

	return players, nil
}

func (m *postgresPlayerRepo) CreatePlayerTransaction(player models.Player, ctx context.Context, tx *sql.Tx) error {
	stmt := `insert into players (league_id, user_id, is_commissioner, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := tx.ExecContext(
		ctx,
		stmt,
		player.LeagueID,
		player.UserID,
		player.IsCommissioner,
		player.IsActive,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
