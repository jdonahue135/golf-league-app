package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type postgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(conn *sql.DB) repository.UserRepo {
	return &postgresUserRepo{
		DB: conn,
	}
}

func (m *postgresUserRepo) AllUsers() bool {
	return true
}

// GetUserByID returns a user by id
func (m *postgresUserRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level_id, created_at, updated_at from users where id=$1`

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
func (m *postgresUserRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, access_level_id, created_at, updated_at from users where email=$1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
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
func (m *postgresUserRepo) UpdateUser(u models.User) error {
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

// Authenticate authenticates a user
func (m *postgresUserRepo) Authenticate(email, password string) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var accessLevel int

	row := m.DB.QueryRowContext(ctx, "select id, access_level_id, password from users where email = $1", email)
	err := row.Scan(&id, &accessLevel, &hashedPassword)
	if err != nil {
		return id, accessLevel, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, 0, err
	}

	return id, accessLevel, nil
}

// CreateUser creates a user
func (m *postgresUserRepo) CreateUser(u models.User, password string) (int, error) {
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

func (m *postgresUserRepo) CreateInactiveUserTransaction(u models.User, ctx context.Context, tx *sql.Tx) (int, error) {
	var userID int

	stmt := `insert into users (first_name, last_name, email, access_level_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) returning id`

	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		time.Now(),
	).Scan(&userID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return userID, nil
}

func (m *postgresUserRepo) BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	tx, err := m.DB.BeginTx(ctx, nil)
	return ctx, cancel, tx, err
}
