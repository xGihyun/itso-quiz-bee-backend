package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (dr *DatabaseRepository) GetResults(ctx context.Context, quizID string) ([]Result, error) {
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

	rows, err := dr.Querier.Query(ctx, sql, quizID)
	if err != nil {
		return nil, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[Result])
	if err != nil {
		return nil, err
	}

	return results, nil
}
