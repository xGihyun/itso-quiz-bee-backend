package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type LiveUpdateStatusRequest struct {
	QuizID string `json:"quiz_id"`
	Status Status `json:"status"`
}

func (r *repository) LiveUpdateStatus(ctx context.Context, data LiveUpdateStatusRequest) (Question, error) {
	tx, err := r.querier.Begin(ctx)
	if err != nil {
		return Question{}, err
	}

	var question Question

	err = database.Transaction(ctx, tx, func() error {
		sql := `
        UPDATE quizzes 
        SET status = ($1)
        WHERE quiz_id = ($2)
        `

		if _, err := tx.Exec(ctx, sql, data.Status, data.QuizID); err != nil {
			return err
		}

		if data.Status != Started {
			return nil
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

		row := r.querier.QueryRow(ctx, sql, data.QuizID)
		if err := row.Scan(
			&question.QuizQuestionID,
			&question.Content,
			&question.Variant,
			&question.Points,
			&question.OrderNumber,
			&question.Answers,
		); err != nil {
			return err
		}

		sql = `
        UPDATE players_in_quizzes
        SET quiz_question_id = ($1)
        WHERE quiz_id = ($2)
        `

		if _, err := tx.Exec(ctx, sql, question.QuizQuestionID, data.QuizID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Question{}, err
	}

	return question, nil
}

type LiveUpdateQuestionRequest struct {
	Question
	QuizID string `json:"quiz_id"`
}

func (r *repository) LiveUpdateQuestion(ctx context.Context, data LiveUpdateQuestionRequest) error {
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

type LiveSubmitAnswerRequest struct {
	CreateWrittenAnswerRequest
}

func (r *repository) LiveSubmitAnswer(ctx context.Context, data LiveSubmitAnswerRequest) error {
	sql := `
	INSERT INTO player_written_answers (content, quiz_question_id, user_id)
	VALUES ($1, $2, $3)
    `

	if _, err := r.querier.Exec(ctx, sql, data.Content, data.QuizQuestionID, data.UserID); err != nil {
		return err
	}

	return nil
}
