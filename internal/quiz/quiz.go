package quiz

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
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

func (r *repository) Save(ctx context.Context, data Quiz) error {
	sql := `
    INSERT INTO quizzes (quiz_id, name, description, status)
    VALUES ($1, $2, $3, $4)
	ON CONFLICT(quiz_id)
	DO UPDATE SET
		name = ($2),
		description = ($3),
		status = ($4)
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

		for i, question := range data.Questions {
			question.OrderNumber = int16(i + 1)
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

func createQuestion(
	querier database.Querier,
	ctx context.Context,
	question Question,
	quizID string,
) error {
	sql := `
    INSERT INTO quiz_questions (
        quiz_question_id, 
        content, 
        points, 
        order_number, 
        duration, 
        quiz_id
    )
    VALUES (
        $1, $2, $3, $4, 
            CASE WHEN $5::int IS NOT NULL 
            THEN make_interval(secs => $5::int)
            ELSE NULL
        END,
        $6
    )
    ON CONFLICT(quiz_question_id)
    DO UPDATE SET
        content = ($2),
        points = ($3),
        order_number = ($4),
        duration = 
            CASE WHEN $5::int IS NOT NULL 
            THEN make_interval(secs => $5::int)
            ELSE NULL
            END
    RETURNING quiz_question_id
    `

	row := querier.QueryRow(
		ctx,
		sql,
		question.QuizQuestionID,
		question.Content,
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

func createAnswer(
	querier database.Querier,
	ctx context.Context,
	answer Answer,
	questionID string,
) error {
	sql := `
	INSERT INTO quiz_answers (quiz_answer_id, content, quiz_question_id)
	VALUES ($1, $2, $3)
	ON CONFLICT(quiz_answer_id)
	DO UPDATE SET content = ($2)
    `

	if _, err := querier.Exec(
		ctx,
		sql,
		answer.QuizAnswerID,
		answer.Content,
		questionID,
	); err != nil {
		return err
	}

	return nil
}

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
