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
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Content-Type", "application/json")

	if _, err := r.Cookie("session"); err != http.ErrNoCookie {
		return api.Response{
			Error:      err,
			Message:    "User session exists.",
			StatusCode: http.StatusConflict,
		}
	}

	ctx := context.Background()

	var data LoginRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
	SELECT user_id, email, role FROM users
	WHERE email = ($1) AND password = ($2)
    `

	row := d.DB.QueryRow(ctx, sql, data.Email, data.Password)

	var user user.User

	if err := row.Scan(&user.UserID, &user.Email, &user.Role); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
		}
	}

	// TODO: Change value to something else
	cookie := http.Cookie{
		Name:     "session",
		Value:    user.UserID,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	return api.Response{}
}
