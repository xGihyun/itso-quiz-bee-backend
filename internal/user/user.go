package user

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type Dependency struct {
	DB *pgxpool.Pool
}

type Role string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type User struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
}

func (d Dependency) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	sql := `
	INSERT INTO users (email, password, role)
	VALUES ($1, $2, $3)
	`

	if _, err := d.DB.Exec(ctx, sql, "gihyun@email.com", "password", Player); err != nil {
		log.Print("Something went wrong: ", err)
		http.Error(w, "Something went wrong", 500)

		return
	}

	log.Print("Create new user!")
	w.WriteHeader(http.StatusCreated)
}

func (d Dependency) GetByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	query := "SELECT user_id, email, role FROM users WHERE user_id = ($1)"

	id := r.PathValue("id")

	row := d.DB.QueryRow(ctx, query, id)

	var user User

	if err := row.Scan(&user.UserID, &user.Email, &user.Role); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
		}
	}

	if err := api.WriteJSON(w, user); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
