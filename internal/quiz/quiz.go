package quiz

import (
	"context"
	"time"
)

type BasicInfo struct {
	QuizID      string  `json:"quiz_id"`
	Name        *string  `json:"name"`
	Description *string `json:"description"`
	Status      *Status  `json:"status"`
	IsTimerAuto *bool    `json:"is_timer_auto"`
}

type Quiz struct {
	QuizID      string     `json:"quiz_id"`
	CreatedAt   time.Time  `json:"created_at"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      Status     `json:"status"`
	Questions   []Question `json:"questions"`
	IsTimerAuto bool       `json:"is_timer_auto"`
}

func (r *repository) UpdateBasicInfo(ctx context.Context, data BasicInfo) error {
	sql := `
    UPDATE quizzes 
    SET
        name = COALESCE($1, name),
        description = COALESCE($2, description),
        status = COALESCE($3, status),
        is_timer_auto = COALESCE($4, is_timer_auto)
    WHERE quiz_id = ($5)
    `

	if _, err := r.querier.Exec(
		ctx,
		sql,
		data.Name,
		data.Description,
		data.Status,
		data.IsTimerAuto,
		data.QuizID,
	); err != nil {
		return err
	}

	return nil
}
