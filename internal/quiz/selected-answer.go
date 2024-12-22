package quiz

import "context"

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
