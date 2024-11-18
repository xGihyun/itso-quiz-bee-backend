package quiz

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type Status string

const (
	Open    Status = "open"
	Started Status = "started"
	Paused  Status = "paused"
	Closed  Status = "closed"
)

type NewQuizRequest struct {
	BasicInfo
	Questions []NewQuestion `json:"questions"`
}

// TODO: Use transactions
func (dr *DatabaseRepository) Create(ctx context.Context, data NewQuizRequest) error {
	sql := `
    INSERT INTO quizzes (quiz_id, name, description, status, lobby_id)
    VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT(quiz_id)
	DO UPDATE SET
		name = ($2),
		description = ($3),
		status = ($4),
		lobby_id = ($5)
    RETURNING quiz_id
    `

	// NOTE: This `tx` won't work
	tx, err := dr.Querier.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	// if data.LobbyID != nil && *data.LobbyID == "" {
	// 	data.LobbyID = nil
	// }

	_, err = dr.Querier.Exec(ctx, sql, data.QuizID, data.Name, data.Description, data.Status, data.LobbyID)
	if err != nil {
		return err
	}

	for i, question := range data.Questions {
		if err := dr.CreateQuestion(ctx, question, data.QuizID, i+1); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type QuizResponse struct {
	BasicInfo
	Questions []Question `json:"questions"`
}

type Answer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	Content      string `json:"content"`
	IsCorrect    bool   `json:"is_correct"`
}

func (dr *DatabaseRepository) GetByID(ctx context.Context, quizID string) (QuizResponse, error) {
	sql := `
	SELECT 
		quizzes.quiz_id, 
		quizzes.name, 
		quizzes.description,
		quizzes.lobby_id,
		quizzes.status,
		(
			SELECT jsonb_agg(
				jsonb_build_object(
					'quiz_question_id', quiz_questions.quiz_question_id,
					'content', quiz_questions.content,
					'variant', quiz_questions.variant,
					'points', quiz_questions.points,
					'order_number', quiz_questions.order_number,
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

	row := dr.Querier.QueryRow(ctx, sql, quizID)

	var quiz QuizResponse

	if err := row.Scan(&quiz.QuizID, &quiz.Name, &quiz.Description, &quiz.LobbyID, &quiz.Status, &quiz.Questions); err != nil {
		return QuizResponse{}, err
	}

	return quiz, nil
}

type BasicInfo struct {
	QuizID      string  `json:"quiz_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      Status  `json:"status"`
	LobbyID     *string `json:"lobby_id"`
}

func (dr *DatabaseRepository) GetAll(ctx context.Context) ([]BasicInfo, error) {
	sql := `
	SELECT quiz_id, name, description, lobby_id, status
	FROM quizzes
	`

	rows, err := dr.Querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	quizzes, err := pgx.CollectRows(rows, pgx.RowToStructByName[BasicInfo])
	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (dr *DatabaseRepository) UpdateByID(ctx context.Context, data BasicInfo) error {
	sql := `
	UPDATE quizzes
	SET name = ($1), description = ($2), lobby_id = ($3), status = ($4)
	WHERE quiz_id = ($5)
	`

	if _, err := dr.Querier.Exec(ctx, sql, data.Name, data.Description, data.LobbyID, data.Status, data.QuizID); err != nil {
		return err
	}

	return nil
}

type UpdateStatusRequest struct {
	QuizID string `json:"quiz_id"`
	Status Status `json:"status"`
}

func (dr *DatabaseRepository) UpdateStatusByID(ctx context.Context, data UpdateStatusRequest) error {
	// NOTE: Testing `squirrel` query builder
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, _, err := psql.Update("quizzes").Set("status", data.Status).Where("quiz_id = (?)", data.QuizID).ToSql()
	if err != nil {
		return err
	}

	if _, err := dr.Querier.Exec(ctx, sql, data.Status, data.QuizID); err != nil {
		return err
	}

	return nil
}
