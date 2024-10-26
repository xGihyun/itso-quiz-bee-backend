package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) HandleCreate(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data UserRequest

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

	return api.Response{StatusCode: http.StatusCreated}
}

func (s *Service) HandleGetByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	id := r.PathValue("user_id")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := api.WriteJSON(w, user); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{}
}
