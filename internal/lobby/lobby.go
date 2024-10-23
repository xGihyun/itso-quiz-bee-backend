package lobby

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type Dependency struct {
	DB *pgxpool.Pool
}

type Status string

const (
	Open   Status = "open"
	Closed Status = "closed"
)

type NewLobby struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Capacity    sql.NullInt16  `json:"capacity"`
	Status      Status         `json:"status"`
}

const OTP_LENGTH = 6

// TODO:
// - Separate some logic in different functions
// - Use transactions
func (d Dependency) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewLobby

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    INSERT INTO lobbies (name, description, capacity, status)
    VALUES ($1, $2, $3, $4)
    RETURNING lobby_id
    `

	row := d.DB.QueryRow(ctx, sql, data.Name, data.Description, data.Capacity, data.Status)

	var lobbyID string

	if err := row.Scan(&lobbyID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	code, err := GenerateOTP(OTP_LENGTH)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	sql = `
    INSERT INTO lobby_codes (code, lobby_id)
    VALUES ($1, $2)
    `

	if _, err := d.DB.Exec(ctx, sql, code, lobbyID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}

type JoinRequestData struct {
	Code   string `json:"code"`
	UserID string `json:"user_id"`
}

func (d Dependency) Join(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data JoinRequestData

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    SELECT lobby_id FROM lobby_codes
    WHERE code = ($1)
    `

	var lobbyID string

	row := d.DB.QueryRow(ctx, sql, data.Code)

	if err := row.Scan(&lobbyID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
			Message:    "Lobby with code " + data.Code + " not found.",
		}
	}

	sql = `
    INSERT INTO users_in_lobbies (user_id, lobby_id)
    VALUES ($1, $2)
    `

	if _, err := d.DB.Exec(ctx, sql, data.UserID, lobbyID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
