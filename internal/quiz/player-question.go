package quiz

import "context"

type UpdatePlayersQuestion struct {
	Question
	QuizID string `json:"quiz_id"`
}

func (r *repository) UpdatePlayersQuestion(ctx context.Context, data UpdatePlayersQuestion) error {
	sql := `
	UPDATE players_in_quizzes
	SET quiz_question_id = ($1)
	WHERE quiz_id = ($2)
	`

	if _, err := r.querier.Exec(ctx, sql, data.QuizQuestionID, data.QuizID); err != nil {
		return err
	}

	return nil
}
