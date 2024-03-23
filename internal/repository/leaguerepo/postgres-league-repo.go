package leaguerepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type postgresLeagueRepo struct {
	DB *sql.DB
}

func NewPostgresLeagueRepo(conn *sql.DB) repository.LeagueRepo {
	return &postgresLeagueRepo{
		DB: conn,
	}
}

// GetLeagueByName returns a league by name
func (m *postgresLeagueRepo) GetLeagueByName(name string) (models.League, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, created_at, updated_at from leagues where name=$1`

	row := m.DB.QueryRowContext(ctx, query, name)

	var l models.League

	err := row.Scan(
		&l.ID,
		&l.Name,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return l, err
	}

	return l, nil
}

// GetLeagueByID returns a league by ID
func (m *postgresLeagueRepo) GetLeagueByID(id int) (models.League, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, created_at, updated_at from leagues where id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var l models.League

	err := row.Scan(
		&l.ID,
		&l.Name,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return l, err
	}

	return l, nil
}

func (m *postgresLeagueRepo) GetLeaguesByUserID(userID int) ([]models.League, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select
	l.id, l.name, l.created_at, l.updated_at 
	from leagues l 
	join players p on l.id = p.league_id
	where p.user_id=$1`

	var leagues []models.League

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return leagues, err
	}

	defer rows.Close()

	for rows.Next() {
		var l models.League

		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return leagues, err
		}

		leagues = append(leagues, l)
	}

	if err = rows.Err(); err != nil {
		return leagues, err
	}

	return leagues, nil
}

func (m *postgresLeagueRepo) CreateLeagueTransaction(league models.League, ctx context.Context, tx *sql.Tx) (int, error) {
	var leagueID int
	stmt := `insert into leagues (name, created_at, updated_at) values ($1, $2, $3) returning id`

	err := tx.QueryRowContext(
		ctx,
		stmt,
		league.Name,
		time.Now(),
		time.Now(),
	).Scan(&leagueID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return leagueID, nil
}

func (m *postgresLeagueRepo) CreateLeague(league models.League, commissioner models.Player) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var leagueID int

	// Begin a transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	stmt := `insert into leagues (name, created_at, updated_at) values ($1, $2, $3) returning id`

	err = tx.QueryRowContext(
		ctx,
		stmt,
		league.Name,
		time.Now(),
		time.Now(),
	).Scan(&leagueID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt = `insert into players (league_id, user_id, is_commissioner, is_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(
		ctx,
		stmt,
		leagueID,
		commissioner.UserID,
		commissioner.IsCommissioner,
		commissioner.IsActive,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction if all operations are successful
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return leagueID, nil
}
