package userrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type testUserRepo struct{}

func NewTestUserRepo() repository.UserRepo {
	return &testUserRepo{}
}

func (m *testUserRepo) CreateUser(u models.User, password string) (int, error) {
	if password == "error" {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testUserRepo) Authenticate(email, password string) (int, int, error) {
	if email == "jack@nimble.com" {
		return 0, 0, errors.New("some error")
	}
	return 1, 1, nil
}

func (m *testUserRepo) AllUsers() bool {
	return true
}

func (m *testUserRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	if email == "me@here.ca" {
		return u, errors.New("some error")
	}
	return u, nil
}

func (m *testUserRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	if id == 0 {
		return u, errors.New("some error")
	}
	return u, nil
}

func (m *testUserRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testUserRepo) CreateInactiveUserTransaction(u models.User, ctx context.Context, tx *sql.Tx) (int, error) {
	if u.FirstName == "user create error" {
		return 1, errors.New("some error")
	}
	if u.FirstName == "player create error" {
		return 2, nil
	}

	return 1, nil
}
