package lobby

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

	var data NewLobbyRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	lobby, err := s.repo.Create(ctx, data)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Fail,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Data: lobby, Status: api.Success, Message: "Created lobby."}
}

func (s *Service) Join(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data JoinRequest

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
				Status:     api.Fail,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
				Status:     api.Error,
			}
		}
	}

	data.UserID = cookie.Value

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	lobby, err := s.repo.Join(ctx, data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
				Message:    "Lobby with code " + data.Code + " not found.",
				Status:     api.Fail,
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success, Message: "Joined lobby.", Data: lobby}
}
