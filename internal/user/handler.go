package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data CreateUserRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Message: "Successfully created user."}
}

func (s *Service) GetByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	id := r.PathValue("user_id")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
				Message:    "User not found.",
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch user.",
		}
	}

	return api.Response{
		Data:       user,
		StatusCode: http.StatusOK,
		Message:    "Fetched user.",
	}
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{
		Data:       users,
		StatusCode: http.StatusOK,
		Message:    "Fetched all users.",
	}
}
