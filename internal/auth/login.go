package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Service) Login(w http.ResponseWriter, r *http.Request) api.Response {
	// TODO: Save session IDs on the server side using a more random and secure value.
	// if _, err := r.Cookie("session"); err != http.ErrNoCookie {
	// 	return api.Response{
	// 		Error:      err,
	// 		StatusCode: http.StatusConflict,
	// 		Message:    "User session already exists.",
	// 	}
	// }
	
	ctx := context.Background()

	var data LoginRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	sql := `
	SELECT user_id, created_at, username, name, role 
    FROM users
	WHERE username = ($1) AND password = ($2)
    `

	row := s.querier.QueryRow(ctx, sql, data.Username, data.Password)

	var user user.GetUserResponse

	if err := row.Scan(&user.UserID, &user.CreatedAt, &user.Username, &user.Name, &user.Role); err != nil {
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

	// TODO: Change `Value` to a more secure value
	// cookie := http.Cookie{
	// 	Name:     "session",
	// 	Value:    user.UserID,
	// 	Path:     "/",
	// 	SameSite: http.SameSiteNoneMode,
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	// Domain:   "http://192.168.1.2:3001",
	// }
	// http.SetCookie(w, &cookie)

	return api.Response{
		StatusCode: http.StatusOK,
		Data:       user,
		Message:    "Successfully logged in.",
	}
}
