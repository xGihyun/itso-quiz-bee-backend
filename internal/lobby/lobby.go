package lobby

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Dependency struct {
	DB database.Querier
}

type Status string

const (
	Open   Status = "open"
	Closed Status = "closed"
)

type NewLobbyRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int16  `json:"capacity"`
	Status      Status  `json:"status"`
}

type NewLobbyResponse struct {
	LobbyID string `json:"lobby_id"`
	Code    string `json:"code"`
}

const OTP_LENGTH = 6

// TODO:
// - Separate some logic in different functions
// - Use transactions
func (d Dependency) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewLobbyRequest

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
		RETURNING code, lobby_id
		`

	row = d.DB.QueryRow(ctx, sql, code, lobbyID)

	var lobby NewLobbyResponse

	if err := row.Scan(&lobby.Code, &lobby.LobbyID); err != nil {
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

type JoinRequest struct {
	Code   string `json:"code"`
	UserID string `json:"user_id"`
}

func (d Dependency) Join(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data JoinRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    SELECT lobby_codes.lobby_id FROM lobby_codes
	JOIN lobbies ON lobbies.lobby_id = lobby_codes.lobby_id
    WHERE lobby_codes.code = ($1) AND lobbies.status = 'open'
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
