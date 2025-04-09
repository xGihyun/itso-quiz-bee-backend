package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	var data createUserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error: fmt.Errorf("create user: %w", err),
			Code:  http.StatusBadRequest,
		}
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return api.Response{
			Error: fmt.Errorf("create user: %w", err),
			Code:  http.StatusInternalServerError,
		}
	}

	return api.Response{Code: http.StatusCreated, Message: "Successfully created user."}
}

func (s *Service) GetByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	id := r.PathValue("user_id")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:   fmt.Errorf("get user by ID: %w", err),
				Code:    http.StatusNotFound,
				Message: "User not found.",
			}
		}

		return api.Response{
			Error:   fmt.Errorf("get user by ID: %w", err),
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch user.",
		}
	}

	return api.Response{
		Data:    user,
		Code:    http.StatusOK,
		Message: "Fetched user.",
	}
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return api.Response{
			Error: fmt.Errorf("get users: %w", err),
			Code:  http.StatusInternalServerError,
		}
	}

	return api.Response{
		Data:    users,
		Code:    http.StatusOK,
		Message: "Fetched all users.",
	}
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

func (s *Service) SignIn(w http.ResponseWriter, r *http.Request) api.Response {
    ctx := r.Context()

	var data signInRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   fmt.Errorf("sign in: %w", err),
			Code:    http.StatusBadRequest,
			Message: "Invalid sign in request.",
		}
	}

	response, err := s.repo.SignIn(ctx, data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:   fmt.Errorf("sign in: %w", err),
				Code:    http.StatusNotFound,
				Message: "Invalid credentials.",
			}
		}

		if errors.Is(err, errInvalidPassword) {
			return api.Response{
				Error:   fmt.Errorf("sign in: %w", err),
				Code:    http.StatusUnauthorized,
				Message: "Invalid password.",
			}
		}

		return api.Response{
			Error:   fmt.Errorf("sign in: %w", err),
			Code:    http.StatusInternalServerError,
			Message: "Failed to sign in.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Message: "Successfully signed in.",
		Data:    response,
	}
}

type signOutRequest struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

func (s *Service) SignOut(w http.ResponseWriter, r *http.Request) api.Response {
    ctx := r.Context()

	var data signOutRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   fmt.Errorf("sign out: %w", err),
			Code:    http.StatusBadRequest,
			Message: "Invalid sign out request.",
		}
	}

	if err := s.repo.invalidateSession(ctx, data.Token, data.UserID); err != nil {
		return api.Response{
			Error:   fmt.Errorf("sign out: %w", err),
			Code:    http.StatusInternalServerError,
			Message: "Failed to sign out.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Message: "Successfully signed out.",
	}
}

func (s *Service) GetSession(w http.ResponseWriter, r *http.Request) api.Response {
    ctx := r.Context()

	token := r.URL.Query().Get("token")

	result, err := s.repo.validateSessionToken(ctx, token)
	if err != nil {
		return api.Response{
			Error:   fmt.Errorf("get session: %w", err),
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user session.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Message: "Successfully fetched user session.",
		Data:    result,
	}
}
