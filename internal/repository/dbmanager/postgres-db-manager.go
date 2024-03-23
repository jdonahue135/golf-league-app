package dbmanager

import (
	"context"
	"database/sql"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/repository"
)

type postgresDBManager struct {
	DB *sql.DB
}

func NewPostgresDBManager(conn *sql.DB) repository.DBManager {
	return &postgresDBManager{
		DB: conn,
	}
}

func (m *postgresDBManager) CommitTransaction(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m *postgresDBManager) BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	tx, err := m.DB.BeginTx(ctx, nil)
	return ctx, cancel, tx, err
}
