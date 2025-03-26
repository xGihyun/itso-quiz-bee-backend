package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Status string

const (
	Open    Status = "open"
	Started Status = "started"
	Paused  Status = "paused"
	Closed  Status = "closed"
)

type UpdateStatusRequest struct {
	QuizID string `json:"quizId"`
	Status Status `json:"status"`
}

func (r *repository) UpdateStatus(ctx context.Context, data UpdateStatusRequest) error {
	sql := `
        UPDATE quizzes 
        SET status = ($1)
        WHERE quiz_id = ($2)
        `

	if _, err := r.querier.Exec(ctx, sql, data.Status, data.QuizID); err != nil {
		return err
	}

	return nil
}

func (r *repository) Start(ctx context.Context, quizID string) (Question, error) {
	tx, err := r.querier.Begin(ctx)
	if err != nil {
		return Question{}, err
	}

	var question Question

	err = database.Transaction(ctx, tx, func() error {
		sql := `
        UPDATE quizzes 
        SET status = 'started'
        WHERE quiz_id = ($1)
        `

		if _, err := tx.Exec(ctx, sql, quizID); err != nil {
			return err
		}

		sql = `
        SELECT
            quiz_question_id,
            content,
            variant,
            points,
            order_number,
            EXTRACT(epoch FROM duration)::INT AS duration,
            (
                SELECT jsonb_agg(
                    jsonb_build_object(
                        'quiz_answer_id', quiz_answer_id,
                        'content', content,
                        'is_correct', is_correct
                    )
                )
                FROM quiz_answers
                WHERE quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
            ) AS answers
        FROM quiz_questions
        WHERE quiz_id = ($1)
        AND order_number = 1
        LIMIT 1
        `

		row := r.querier.QueryRow(ctx, sql, quizID)
		if err := row.Scan(
			&question.QuizQuestionID,
			&question.Content,
			&question.Variant,
			&question.Points,
			&question.OrderNumber,
			&question.Duration,
			&question.Answers,
		); err != nil {
			return err
		}

		sql = `
        UPDATE players_in_quizzes
        SET quiz_question_id = ($1)
        WHERE quiz_id = ($2)
        `

		if _, err := tx.Exec(ctx, sql, question.QuizQuestionID, quizID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Question{}, err
	}

	return question, nil
}
