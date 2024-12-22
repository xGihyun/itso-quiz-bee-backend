package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type AddPlayerRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (r *repository) AddPlayer(ctx context.Context, data AddPlayerRequest) error {
	sql := `
	INSERT INTO players_in_quizzes (user_id, quiz_id)
	VALUES ($1, $2)
	ON CONFLICT(user_id, quiz_id)
	DO NOTHING
	`

	if _, err := r.querier.Exec(ctx, sql, data.UserID, data.QuizID); err != nil {
		return err
	}

	return nil
}

type UpdatePlayerQuestionRequest struct {
	QuizID         string `json:"quiz_id"`
	QuizQuestionID string `json:"quiz_question_id"`
}

type UpdatePlayerQuestionResponse struct {
	QuizQuestionID string `json:"quiz_question_id"`
}

func (r *repository) UpdatePlayerQuestion(ctx context.Context, data UpdatePlayerQuestionRequest) error {
	sql := `
	UPDATE players_in_quizzes
	SET quiz_question_id = ($1)
	WHERE quiz_id = ($2)
	`

	if _, err := r.querier.Exec(ctx, sql, data.QuizQuestionID, data.QuizID); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetPlayers(ctx context.Context, quizID string) ([]user.GetUserResponse, error) {
	sql := `
	SELECT 
		users.user_id,
		users.created_at,
		users.username,
		users.role,
		users.name
	FROM players_in_quizzes
	JOIN users ON users.user_id = players_in_quizzes.user_id
	WHERE players_in_quizzes.quiz_id = ($1)
	`

	rows, err := r.querier.Query(ctx, sql, quizID)
	if err != nil {
		return []user.GetUserResponse{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.GetUserResponse])
	if err != nil {
		return []user.GetUserResponse{}, err
	}

	return users, nil
}
