package quiz

import "context"

type UpdatePlayersQuestionRequest struct {
	Question Question `json:"question"`
	QuizID   string   `json:"quiz_id"`
}

func (r *repository) UpdatePlayersQuestion(
	ctx context.Context,
	data UpdatePlayersQuestionRequest,
) error {
	sql := `
	UPDATE players_in_quizzes
	SET quiz_question_id = ($1)
	WHERE quiz_id = ($2)
	`

	if _, err := r.querier.Exec(ctx, sql, data.Question.QuizQuestionID, data.QuizID); err != nil {
		return err
	}

	return nil
}
