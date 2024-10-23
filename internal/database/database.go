package database

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
