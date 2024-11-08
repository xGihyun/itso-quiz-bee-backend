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

	var data UserRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success}
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
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{
		Data:       user,
		StatusCode: http.StatusOK,
		Status:     api.Success,
	}
}
