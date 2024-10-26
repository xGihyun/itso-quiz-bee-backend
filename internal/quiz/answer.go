package quiz

import (
	"context"
)

type NewAnswer struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

func (dr *DatabaseRepository) CreateAnswer(ctx context.Context, answer NewAnswer, questionID string) error {
	sql := `
	INSERT INTO quiz_answers (content, is_correct, quiz_question_id)
	VALUES ($1, $2, $3)
    `

	if _, err := dr.Querier.Exec(ctx, sql, answer.Content, answer.IsCorrect, questionID); err != nil {
		return err
	}

	return nil
}

type NewSelectedAnswer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	UserID       string `json:"user_id"`
}

func (dr *DatabaseRepository) CreateSelectedAnswer(ctx context.Context, data NewSelectedAnswer) error {
	sql := `
	INSERT INTO player_selected_answers (quiz_answer_id, user_id)
	VALUES ($1, $2)
    `

	if _, err := dr.Querier.Exec(ctx, sql, data.QuizAnswerID, data.UserID); err != nil {
		return err
	}

	return nil
}
