package quiz

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type Result struct {
	Score  int16  `json:"score"`
	UserID string `json:"user_id"`
	// Answers []PlayerAnswer `json:"answers"`
}

type PlayerAnswer struct {
	PlayerSelectedAnswerID string `json:"player_selected_answer_id"`
	IsCorrect              bool   `json:"is_correct"`
}

func (d Dependency) GetResults(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	// TODO:
	// - Also get score from written answers and add it with selected answers
	// - Put SQL queries in their own .sql files maybe (?)

	sql := `
	WITH player_scores AS (
		SELECT 
			SUM(quiz_questions.points) AS score,
			player_selected_answers.player_selected_answer_id,
			player_selected_answers.user_id
		FROM player_selected_answers
		JOIN quiz_answers ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE quiz_answers.is_correct IS TRUE AND quiz_questions.quiz_id = ($1)
		GROUP BY 
			player_selected_answers.user_id, 
			player_selected_answers.player_selected_answer_id
	),
	player_answers AS (
		SELECT 
			player_selected_answers.player_selected_answer_id,
			player_selected_answers.user_id,
			quiz_answers.is_correct,
			quiz_questions.order_number,
			quiz_questions.points
		FROM player_selected_answers
		JOIN quiz_answers ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE quiz_questions.quiz_id = ($1)
		ORDER BY quiz_questions.order_number
	)
	SELECT 
		score, 
		user_id
	FROM player_scores
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

	if err := api.WriteJSON(w, results); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{}
}

// func (d Dependency) GetPlayerScore
