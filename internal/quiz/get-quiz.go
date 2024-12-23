package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetByID(ctx context.Context, quizID string) (Quiz, error) {
	sql := `
	SELECT 
		quizzes.quiz_id, 
		quizzes.created_at, 
		quizzes.name, 
		quizzes.description,
		quizzes.status,
		(
			SELECT jsonb_agg(
				jsonb_build_object(
					'quiz_question_id', quiz_questions.quiz_question_id,
					'content', quiz_questions.content,
					'variant', quiz_questions.variant,
					'points', quiz_questions.points,
					'order_number', quiz_questions.order_number,
					'duration', EXTRACT(epoch FROM quiz_questions.duration)::INT,
					'answers', (
						SELECT jsonb_agg(
							jsonb_build_object(
								'quiz_answer_id', quiz_answers.quiz_answer_id,
								'content', quiz_answers.content,
								'is_correct', quiz_answers.is_correct
							)
						)
						FROM quiz_answers
						WHERE quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
					)
				)
			)
			FROM quiz_questions
			WHERE quiz_questions.quiz_id = quizzes.quiz_id
		) as questions
	FROM quizzes
	WHERE quizzes.quiz_id = ($1)
	`

	row := r.querier.QueryRow(ctx, sql, quizID)

	var quiz Quiz

	if err := row.Scan(&quiz.QuizID, &quiz.CreatedAt, &quiz.Name, &quiz.Description, &quiz.Status, &quiz.Questions); err != nil {
		return Quiz{}, err
	}

	return quiz, nil
}

func (r *repository) GetMany(ctx context.Context) ([]BasicInfo, error) {
	sql := `
	SELECT quiz_id, name, description, status
	FROM quizzes
	`

	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	quizzes, err := pgx.CollectRows(rows, pgx.RowToStructByName[BasicInfo])
	if err != nil {
		return []BasicInfo{}, err
	}

	return quizzes, nil
}
