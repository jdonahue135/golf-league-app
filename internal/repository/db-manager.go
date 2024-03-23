package repository

import (
	"context"
	"database/sql"
)

type DBManager interface {
	BeginTransaction() (context.Context, context.CancelFunc, *sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
}
