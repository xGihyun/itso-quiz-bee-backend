package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Dependency struct {
	DB database.Querier
}

type RegisterRequest struct {
	user.UserRequest
	user.Detail
}

func (d Dependency) Register(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data RegisterRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	tx, err := d.DB.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
			Message:    "Failed to start transaction.",
		}
	}

	sql := `
    INSERT INTO users (email, password, role)
    VALUES ($1, $2, $3)
	RETURNING user_id
    `

	row := tx.QueryRow(ctx, sql, data.Email, data.Password, data.Role)

	var userID string

	if err := row.Scan(&userID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	// if err != nil {
	// 	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
	// 		return api.Response{
	// 			Error:      err,
	// 			StatusCode: http.StatusConflict,
	// 			Message:    "User " + data.Email + " already exists.",
	// 			Status:     api.Fail,
	// 		}
	// 	}
	//
	// 	return api.Response{
	// 		Error:      err,
	// 		StatusCode: http.StatusInternalServerError,
	// 		Status:     api.Error,
	// 	}
	// }

	sql = `
	INSERT INTO user_details (user_id, first_name, middle_name, last_name)
	VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.Exec(ctx, sql, userID, data.FirstName, data.MiddleName, data.LastName); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusConflict,
				Message:    "Details of " + data.Email + " already exists.",
				Status:     api.Fail,
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
			Message: "Failed to commit transaction.",
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success, Message: "Succesfully registered."}
}
