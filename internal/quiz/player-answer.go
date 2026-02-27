package quiz

import (
	"context"
)

type PlayerAnswer struct {
	PlayerAnswerID string `json:"playerAnswerid"`
	Content        string `json:"content"`
	IsCorrect      bool   `json:"isCorrect"`
	QuizQuestionID string `json:"quizQuestionId"`
}

type CreateWrittenAnswerRequest struct {
	Content        string `json:"content"`
	QuizQuestionID string `json:"quizQuestionId"`
	UserID         string `json:"userId"`
	QuizID         string `json:"quizId"`
}

func (r *repository) CreateWrittenAnswer(
	ctx context.Context,
	data CreateWrittenAnswerRequest,
) error {
	sql := `
	INSERT INTO player_written_answers (content, quiz_question_id, user_id)
	VALUES ($1, $2, $3)
    RETURNING player_written_answer_id
    `

	if _, err := r.querier.Exec(ctx, sql, data.Content, data.QuizQuestionID, data.UserID); err != nil {
		return err
	}

	return nil
}

type GetWrittenAnswerRequest struct {
	QuizID string `json:"quizId"`
	UserID string `json:"userId"`
}

type GetWrittenAnswerResponse struct {
	PlayerWrittenAnswerID string `json:"playerWrittenAnswerId"`
	Content               string `json:"content"`
}

func (r *repository) GetWrittenAnswer(
	ctx context.Context,
	data GetWrittenAnswerRequest,
) (GetWrittenAnswerResponse, error) {
	question, err := r.GetCurrentQuestion(ctx, data.QuizID)
	if err != nil {
		return GetWrittenAnswerResponse{}, err
	}

	sql := `
	SELECT player_written_answer_id, content
	FROM player_written_answers
	WHERE user_id = ($1) AND quiz_question_id = ($2)
	`

	row := r.querier.QueryRow(ctx, sql, data.UserID, question.Question.QuizQuestionID)

	var answer GetWrittenAnswerResponse

	if err := row.Scan(&answer.PlayerWrittenAnswerID, &answer.Content); err != nil {
		return GetWrittenAnswerResponse{}, err
	}

	return answer, nil
}
