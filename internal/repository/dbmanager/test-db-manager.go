package dbmanager

import (
	"context"
	"database/sql"

	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type testDBManager struct{}

func NewTestDBManager() repository.DBManager {
	return &testDBManager{}
}

func (m *testDBManager) BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error) {
	return context.Background(), func() {}, &sql.Tx{}, nil
}

func (m *testDBManager) CommitTransaction(tx *sql.Tx) error {
	return nil
}
