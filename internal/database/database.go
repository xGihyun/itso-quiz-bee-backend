package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Postgres struct {
	DB *pgxpool.Pool
}

func WithTx(pg Postgres) error {
	return nil
}

// import (
// 	"net/http"
//
// 	"github.com/jackc/pgx/v5"
// 	"github.com/xGihyun/itso-quiz-bee/internal/api"
// )
//
// func TxCommitOrRollback(tx pgx.Tx, res api.Response) api.Response {
// 	if res.Error != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return api.Response{
// 				Error:      err,
// 				StatusCode: http.StatusInternalServerError,
// 			}
// 		}
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return api.Response{
// 			Error:      err,
// 			StatusCode: http.StatusInternalServerError,
// 		}
// 	}
//
// 	return api.Response{}
// }
