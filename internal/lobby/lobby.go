package lobby

import (
	"context"

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
// - Use transactions
func (dr *DatabaseRepository) Create(ctx context.Context, data NewLobbyRequest) (NewLobbyResponse, error) {
	sql := `
		INSERT INTO lobbies (name, description, capacity, status)
		VALUES ($1, $2, $3, $4)
		RETURNING lobby_id
		`

	row := dr.Querier.QueryRow(ctx, sql, data.Name, data.Description, data.Capacity, data.Status)

	var lobbyID string

	if err := row.Scan(&lobbyID); err != nil {
		return NewLobbyResponse{}, err
	}

	code, err := GenerateOTP(OTP_LENGTH)
	if err != nil {
		return NewLobbyResponse{}, err
	}

	sql = `
		INSERT INTO lobby_codes (code, lobby_id)
	    VALUES ($1, $2)
		RETURNING code, lobby_id
		`

	row = dr.Querier.QueryRow(ctx, sql, code, lobbyID)

	var lobby NewLobbyResponse

	if err := row.Scan(&lobby.Code, &lobby.LobbyID); err != nil {
		return NewLobbyResponse{}, err
	}

	return lobby, nil
}

type JoinRequest struct {
	Code   string `json:"code"`
	UserID string `json:"user_id"`
}

type JoinResponse struct {
	LobbyID string `json:"lobby_id"`
}

func (dr *DatabaseRepository) Join(ctx context.Context, data JoinRequest) (JoinResponse, error) {
	sql := `
    SELECT lobby_codes.lobby_id FROM lobby_codes
	JOIN lobbies ON lobbies.lobby_id = lobby_codes.lobby_id
    WHERE lobby_codes.code = ($1) AND lobbies.status = 'open'
    `

	var lobby JoinResponse

	row := dr.Querier.QueryRow(ctx, sql, data.Code)

	if err := row.Scan(&lobby.LobbyID); err != nil {
		return JoinResponse{}, err
	}

	sql = `
    INSERT INTO users_in_lobbies (user_id, lobby_id)
    VALUES ($1, $2)
    `

	if _, err := dr.Querier.Exec(ctx, sql, data.UserID, lobby.LobbyID); err != nil {
		return JoinResponse{}, err
	}

	return lobby, nil
}
