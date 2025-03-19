package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

// NOTE: Not sure if this is still needed.
func (s Service) GetCurrentUser(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:   err,
				Message: "User not authenticated.",
				Code:    http.StatusBadRequest,
			}
		default:
			return api.Response{
				Error:   err,
				Message: "Server cookie error.",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	userRepo := user.NewRepository(s.querier)

	user, err := userRepo.GetByID(ctx, cookie.Value)
	if err != nil {
		return api.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch user.",
		}
	}

	return api.Response{
		Code:    http.StatusCreated,
		Message: "Fetched user.",
		Data:    user,
	}
}
