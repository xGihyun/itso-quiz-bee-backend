package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

func (s Service) Register(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data user.CreateUserRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	userRepo := user.NewRepository(s.querier)

	if err := userRepo.Create(ctx, data); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusConflict,
				Message:    "User " + data.Username + " already exists.",
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to register user.",
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Message: "Succesfully registered."}
}
