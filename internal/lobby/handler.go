package lobby

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

	var data NewLobbyRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	lobby, err := s.repo.Create(ctx, data)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := api.WriteJSON(w, lobby); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}

func (s *Service) HandleJoin(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data JoinRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := s.repo.Join(ctx, data); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
				Message:    "Lobby with code " + data.Code + " not found.",
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
