package quiz

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type BasicInfo struct {
	QuizID      string    `json:"quizId"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
}

type Quiz struct {
	BasicInfo
	Questions []Question `json:"questions"`
}

func (r *repository) Get(ctx context.Context, quizID string, includeAnswers bool) (Quiz, error) {
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
					'quizQuestionId', quiz_questions.quiz_question_id,
					'content', quiz_questions.content,
					'variant', quiz_questions.variant,
					'points', quiz_questions.points,
					'orderNumber', quiz_questions.order_number,
					'duration', EXTRACT(epoch FROM quiz_questions.duration)::INT
					%s
				)
			)
			FROM quiz_questions
			WHERE quiz_questions.quiz_id = quizzes.quiz_id
		) as questions
	FROM quizzes
	WHERE quizzes.quiz_id = ($1)
	`

	var answersSql string
	if includeAnswers {
		answersSql = `
		,
		'answers', (
			SELECT jsonb_agg(
				jsonb_build_object(
					'quizAnswerId', quiz_answers.quiz_answer_id,
					'content', quiz_answers.content
				)
			)
			FROM quiz_answers
			WHERE quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
		)
		`
	}

	sql = fmt.Sprintf(sql, answersSql)

	var quiz Quiz

	row := r.querier.QueryRow(ctx, sql, quizID)
	if err := row.Scan(
		&quiz.QuizID,
		&quiz.CreatedAt,
		&quiz.Name,
		&quiz.Description,
		&quiz.Status,
		&quiz.Questions,
	); err != nil {
		return Quiz{}, err
	}

	return quiz, nil
}

func (r *repository) GetBasicInfo(ctx context.Context, quizID string) (BasicInfo, error) {
	sql := `
	SELECT quiz_id, created_at, name, description, status
	FROM quizzes
	WHERE quiz_id = ($1)
	`

	var quiz BasicInfo

	row := r.querier.QueryRow(ctx, sql, quizID)
	if err := row.Scan(
		&quiz.QuizID,
		&quiz.CreatedAt,
		&quiz.Name,
		&quiz.Description,
		&quiz.Status,
	); err != nil {
		return BasicInfo{}, nil
	}

	return quiz, nil
}

func (r *repository) ListBasicInfo(ctx context.Context) ([]BasicInfo, error) {
	sql := `
	SELECT quiz_id, created_at, name, description, status
	FROM quizzes
	`

	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return []BasicInfo{}, err
	}

	quizzes, err := pgx.CollectRows(rows, pgx.RowToStructByName[BasicInfo])
	if err != nil {
		return []BasicInfo{}, err
	}

	return quizzes, nil
}
