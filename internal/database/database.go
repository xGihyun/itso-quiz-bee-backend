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

// type PostgresRepository struct {
// 	Querier Querier
// }
//
// func NewPostgresQuerier(q Querier) *PostgresRepository {
// 	return &PostgresRepository{
// 		Querier: q,
// 	}
// }

// import k
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
