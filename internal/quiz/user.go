package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
)

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

type UpdateCurrentQuestionRequest struct {
	QuizID         string `json:"quiz_id"`
	QuizQuestionID string `json:"quiz_question_id"`
}

type UpdateCurrentQuestionResponse struct {
	QuizQuestionID string `json:"quiz_question_id"`
}

func (dr *DatabaseRepository) UpdateCurrentQuestion(ctx context.Context, data UpdateCurrentQuestionRequest) error {
	sql := `
	UPDATE users_in_quizzes
	SET quiz_question_id = ($1)
	WHERE quiz_id = ($2)
	`

	if _, err := dr.Querier.Exec(ctx, sql, data.QuizQuestionID, data.QuizID); err != nil {
		return err
	}

	return nil
}

type User struct {
	UserID     string  `json:"user_id"`
	FirstName  string  `json:"first_name"`
	MiddleName *string `json:"middle_name"`
	LastName   string  `json:"last_name"`
}

func (dr *DatabaseRepository) GetAllUsers(ctx context.Context, quizID string) ([]User, error) {
	sql := `
	SELECT 
		users_in_quizzes.user_id,
		user_details.first_name,
		user_details.middle_name,
		user_details.last_name
	FROM users_in_quizzes
	JOIN user_details ON user_details.user_id = users_in_quizzes.user_id
	WHERE users_in_quizzes.quiz_id = ($1)
	`

	rows, err := dr.Querier.Query(ctx, sql, quizID)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return users, nil
}
