package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d Dependency) Login(w http.ResponseWriter, r *http.Request) api.Response {
	if _, err := r.Cookie("session"); err != http.ErrNoCookie {
		return api.Response{
			Error:      err,
			Message:    "User session exists.",
			StatusCode: http.StatusConflict,
			Status:     api.Fail,
		}
	}

	ctx := context.Background()

	var data LoginRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	sql := `
	SELECT user_id, email, role FROM users
	WHERE email = ($1) AND password = ($2)
    `

	row := d.DB.QueryRow(ctx, sql, data.Email, data.Password)

	var user user.UserResponse

	if err := row.Scan(&user.UserID, &user.Email, &user.Role); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
			Status:     api.Fail,
		}
	}

	// TODO: Change `Value` to something else
	cookie := http.Cookie{
		Name:     "session",
		Value:    user.UserID,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
		HttpOnly: true,
		// Domain:   "http://192.168.1.2:3001",
	}
	http.SetCookie(w, &cookie)

	return api.Response{StatusCode: http.StatusOK, Status: api.Success, Message: "Successfully logged in.", Data: user}
}
