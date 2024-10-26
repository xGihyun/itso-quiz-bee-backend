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
	// DB *pgxpool.Pool
	DB database.Querier
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d Dependency) Register(w http.ResponseWriter, r *http.Request) api.Response {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	var data RegisterRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    INSERT INTO users (email, password, role)
    VALUES ($1, $2, $3)
    `

	if _, err := d.DB.Exec(ctx, sql, data.Email, data.Password, user.Player); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusConflict,
				Message:    "User " + data.Email + " already exists.",
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
