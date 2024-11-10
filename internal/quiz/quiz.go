package quiz

import (
	"context"
)

type Status string

const (
	Open    Status = "open"
	Started Status = "started"
	Paused  Status = "paused"
	Closed  Status = "closed"
)

type NewQuizRequest struct {
	QuizID      string        `json:"quiz_id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Status      Status        `json:"status"`
	LobbyID     *string       `json:"lobby_id"`
	Questions   []NewQuestion `json:"questions"`
}

// TODO: Use transactions
func (dr *DatabaseRepository) Create(ctx context.Context, data NewQuizRequest) error {
	sql := `
    INSERT INTO quizzes (quiz_id, name, description, status, lobby_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING quiz_id
    `

	// NOTE: This `tx` won't work
	tx, err := dr.Querier.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	if *data.LobbyID == "" {
		data.LobbyID = nil
	}

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

type JoinRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (dr *DatabaseRepository) Join(ctx context.Context, data JoinRequest) error {
	sql := `
	INSERT INTO users_in_quizzes (user_id, quiz_id)
	VALUES ($1, $2)
	`

	if _, err := dr.Querier.Exec(ctx, sql, data.UserID, data.QuizID); err != nil {
		return err
	}

	return nil
}

type NewQuizResponse struct {
	QuizID      string        `json:"quiz_id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Status      Status        `json:"status"`
	LobbyID     *string       `json:"lobby_id"`
	Questions   []NewQuestion `json:"questions"`
}

func (dr *DatabaseRepository) GetByID(ctx context.Context, quizID string) (NewQuizResponse, error) {
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
					'content', quiz_questions.content,
					'variant', quiz_questions.variant,
					'points', quiz_questions.points,
					'answers', (
						SELECT jsonb_agg(
							jsonb_build_object(
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

	var quiz NewQuizResponse

	if err := row.Scan(&quiz.QuizID, &quiz.Name, &quiz.Description, &quiz.LobbyID, &quiz.Status, &quiz.Questions); err != nil {
		return NewQuizResponse{}, err
	}

	return quiz, nil
}
