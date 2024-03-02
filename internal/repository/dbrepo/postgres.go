package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// GetUserByID returns a user by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

// GetUserByEmail returns a user by email
func (m *postgresDBRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where email=$1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the db
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// GetLeagueByName returns a league by name
func (m *postgresDBRepo) GetLeagueByName(name string) (models.League, error) {
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
func (m *postgresDBRepo) GetLeagueByID(id int) (models.League, error) {
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

func (m *postgresDBRepo) CreateLeague(league models.League) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into leagues (name, created_at, updated_at) values ($1, $2, $3) returning id`

	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		league.Name,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) CreateUser(u models.User, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return id, err
	}

	stmt := `insert into users (first_name, last_name, email, password, access_level_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = m.DB.QueryRowContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		string(hashedPassword),
		models.AccessLevelPlayer,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
