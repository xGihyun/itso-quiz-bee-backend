package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

func (r *repository) LiveAddPlayer(ctx context.Context, data AddPlayerRequest) (user.GetUserResponse, error) {
	tx, err := r.querier.Begin(ctx)
	if err != nil {
		return user.GetUserResponse{}, err
	}

	var u user.GetUserResponse

	err = database.Transaction(ctx, tx, func() error {
		sql := `
        INSERT INTO players_in_quizzes (user_id, quiz_id)
        VALUES ($1, $2)
        ON CONFLICT(user_id, quiz_id)
        DO NOTHING
        `

		if _, err := tx.Exec(ctx, sql, data.UserID, data.QuizID); err != nil {
			return err
		}

		sql = `
        SELECT 
            user_id, 
            created_at,
            username,
            role,
            name
        FROM users WHERE user_id = ($1)
        `

		row := r.querier.QueryRow(ctx, sql, data.UserID)

		if err := row.Scan(&u.UserID, &u.CreatedAt, &u.Username, &u.Role, &u.Name); err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return user.GetUserResponse{}, nil
	}

	return u, nil
}
