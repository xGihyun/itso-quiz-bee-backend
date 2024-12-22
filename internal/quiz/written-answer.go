package quiz

import (
	"context"
)

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
