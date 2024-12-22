package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Transaction(ctx context.Context, tx pgx.Tx, fn func() error) error {
	if err := fn(); err != nil {
		_ = tx.Rollback(ctx)

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
