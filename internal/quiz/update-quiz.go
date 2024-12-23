package quiz

import "context"

func (r *repository) UpdateBasicInfo(ctx context.Context, data BasicInfo) error {
	sql := `
    UPDATE quizzes 
    SET
        name = COALESCE($1, name),
        description = COALESCE($2, description),
        status = COALESCE($3, status)
    WHERE quiz_id = ($4)
    `

	if _, err := r.querier.Exec(
		ctx,
		sql,
		data.Name,
		data.Description,
		data.Status,
		data.QuizID,
	); err != nil {
		return err
	}

	return nil
}
