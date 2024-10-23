package quiz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type Result struct {
	Score  int16  `json:"score"`
	UserID string `json:"user_id"`
}

func (d Dependency) GetResults(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	// TODO:
	// - Also get score from written answers and add it with selected answers
	// - Put SQL queries in their own .sql files

	sql := `
	WITH correct_answers AS (
		SELECT 
			player_selected_answers.player_selected_answer_id,
			player_selected_answers.user_id
		FROM player_selected_answers
		JOIN quiz_answers ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE quiz_answers.is_correct IS TRUE AND quiz_questions.quiz_id = ($1)
	)
	SELECT COUNT(*) AS score, user_id FROM correct_answers
	GROUP BY user_id
	`

	rows, err := d.DB.Query(ctx, sql, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[Result])
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{}
}
