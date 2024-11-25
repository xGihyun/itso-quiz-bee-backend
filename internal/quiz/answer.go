package quiz

import (
	"context"
)

type NewAnswer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	Content      string `json:"content"`
	IsCorrect    bool   `json:"is_correct"`
}

func (dr *DatabaseRepository) CreateAnswer(ctx context.Context, answer NewAnswer, questionID string) error {
	sql := `
	INSERT INTO quiz_answers (quiz_answer_id, content, is_correct, quiz_question_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(quiz_answer_id)
	DO UPDATE SET
		content = ($2)
    `

	if _, err := dr.Querier.Exec(ctx, sql, answer.QuizAnswerID, answer.Content, answer.IsCorrect, questionID); err != nil {
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

type NewWrittenAnswerRequest struct {
	Content        string `json:"content"`
	QuizQuestionID string `json:"quiz_question_id"`
	UserID         string `json:"user_id"`
}

func (dr *DatabaseRepository) CreateWrittenAnswer(ctx context.Context, data NewWrittenAnswerRequest) error {
	sql := `
	INSERT INTO player_written_answers (content, quiz_question_id, user_id)
	VALUES ($1, $2, $3)
    `

	if _, err := dr.Querier.Exec(ctx, sql, data.Content, data.QuizQuestionID, data.UserID); err != nil {
		return err
	}

	return nil
}
