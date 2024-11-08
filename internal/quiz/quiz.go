package quiz

import (
	"context"
)

type Status string

const (
	Open    Status = "open"
	Started Status = "started"
	Paused  Status = "paused"
	Closed  Status = "closed"
)

type NewQuiz struct {
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Status      Status        `json:"status"`
	LobbyID     string        `json:"lobby_id"`
	Questions   []NewQuestion `json:"questions"`
}

// TODO: Use transactions
func (dr *DatabaseRepository) Create(ctx context.Context, data NewQuiz) error {
	sql := `
    INSERT INTO quizzes (name, description, status, lobby_id)
    VALUES ($1, $2, $3, $4)
    RETURNING quiz_id
    `

	// NOTE: This `tx` won't work
	tx, err := dr.Querier.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	row := dr.Querier.QueryRow(ctx, sql, data.Name, data.Description, data.Status, data.LobbyID)

	var quizID string

	if err := row.Scan(&quizID); err != nil {
		return err
	}

	for _, question := range data.Questions {
		if err := dr.CreateQuestion(ctx, question, quizID); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type JoinRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (dr *DatabaseRepository) Join(ctx context.Context, data JoinRequest) error {
	sql := `
	INSERT INTO users_in_quizzes (user_id, quiz_id)
	VALUES ($1, $2)
	`

	if _, err := dr.Querier.Exec(ctx, sql, data.UserID, data.QuizID); err != nil {
		return err
	}

	return nil
}
