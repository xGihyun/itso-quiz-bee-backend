package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

func (r *repository) Create(ctx context.Context, data Quiz) error {
	sql := `
    INSERT INTO quizzes (quiz_id, name, description, status)
    VALUES ($1, $2, $3, $4)
	ON CONFLICT(quiz_id)
	DO UPDATE SET
		name = ($2),
		description = ($3),
		status = ($4)
    RETURNING quiz_id
    `

	tx, err := r.querier.Begin(ctx)
	if err != nil {
		return err
	}

	err = database.Transaction(ctx, tx, func() error {
		_, err = tx.Exec(ctx, sql, data.QuizID, data.Name, data.Description, data.Status)
		if err != nil {
			return err
		}

		for _, question := range data.Questions {
			if err := createQuestion(tx, ctx, question, data.QuizID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func createQuestion(querier database.Querier, ctx context.Context, question Question, quizID string) error {
	sql := `
	INSERT INTO quiz_questions (
        quiz_question_id, 
        content, 
        variant, 
        points, 
        order_number, 
        duration, 
        quiz_id
    )
	VALUES (
            $1, $2, $3, $4, $5, 
                CASE WHEN $6 IS NOT NULL
                THEN make_interval(secs => $6)
                ELSE NULL
                END
            $7
        )
	ON CONFLICT(quiz_question_id)
	DO UPDATE SET
		content = ($2),
		variant = ($3),
		points = ($4),
		order_number = ($5),
		duration = 
			CASE WHEN $6 IS NOT NULL
			THEN make_interval(secs => $6)
			ELSE NULL
			END
	RETURNING quiz_question_id
	`

	row := querier.QueryRow(
		ctx,
		sql,
		question.QuizQuestionID,
		question.Content,
		question.Variant,
		question.Points,
		question.OrderNumber,
		question.Duration,
		quizID,
	)

	var questionID string

	if err := row.Scan(&questionID); err != nil {
		return err
	}

	for _, answer := range question.Answers {
		if err := createAnswer(querier, ctx, answer, questionID); err != nil {
			return err
		}
	}

	return nil
}

func createAnswer(querier database.Querier, ctx context.Context, answer Answer, questionID string) error {
	sql := `
	INSERT INTO quiz_answers (quiz_answer_id, content, is_correct, quiz_question_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(quiz_answer_id)
	DO UPDATE SET content = ($2)
    `

	if _, err := querier.Exec(
		ctx,
		sql,
		answer.QuizAnswerID,
		answer.Content,
		answer.IsCorrect,
		questionID,
	); err != nil {
		return err
	}

	return nil
}
