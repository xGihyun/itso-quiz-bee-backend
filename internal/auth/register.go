package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type Model struct {
	DB *pgxpool.Pool
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m Model) Register(w http.ResponseWriter, r *http.Request) api.Response {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	var data RegisterData

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    INSERT INTO users (email, password, role)
    VALUES ($1, $2, $3)
    `

	if _, err := m.DB.Exec(ctx, sql, data.Email, data.Password, "player"); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusConflict,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
