package quiz

import (
	"context"
)

type CreateWrittenAnswerRequest struct {
	Content        string `json:"content"`
	QuizQuestionID string `json:"quiz_question_id"`
	UserID         string `json:"user_id"`
}

func (r *repository) CreateWrittenAnswer(ctx context.Context, data CreateWrittenAnswerRequest) error {
	sql := `
	INSERT INTO player_written_answers (content, quiz_question_id, user_id)
	VALUES ($1, $2, $3)
    `

	if _, err := r.querier.Exec(ctx, sql, data.Content, data.QuizQuestionID, data.UserID); err != nil {
		return err
	}

	return nil
}

type GetWrittenAnswerResponse struct {
	PlayerWrittenAnswerID string `json:"player_written_answer_id"`
	Content               string `json:"content"`
}

func (r *repository) GetWrittenAnswer(ctx context.Context, quizID string, userID string) (GetWrittenAnswerResponse, error) {
	question, err := r.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		return GetWrittenAnswerResponse{}, err
	}

	sql := `
	SELECT player_written_answer_id, content
	FROM player_written_answers
	WHERE user_id = ($1) AND quiz_question_id = ($2)
	`

	row := r.querier.QueryRow(ctx, sql, userID, question.QuizQuestionID)

	var answer GetWrittenAnswerResponse

	if err := row.Scan(&answer.PlayerWrittenAnswerID, &answer.Content); err != nil {
		return GetWrittenAnswerResponse{}, err
	}

	return answer, nil
}

type CreateSelectedAnswerRequest struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	UserID       string `json:"user_id"`
}

func (r *repository) CreateSelectedAnswer(ctx context.Context, data CreateSelectedAnswerRequest) error {
	sql := `
	INSERT INTO player_selected_answers (quiz_answer_id, user_id)
	VALUES ($1, $2)
    `

	if _, err := r.querier.Exec(ctx, sql, data.QuizAnswerID, data.UserID); err != nil {
		return err
	}

	return nil
}
