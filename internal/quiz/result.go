package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Result struct {
	PlayerScore
	Answers []PlayerAnswer `json:"answers"`
}

type PlayerScore struct {
	Score  int16  `json:"score"`
	UserID string `json:"user_id"`
}

type PlayerAnswer struct {
	PlayerSelectedAnswerID string `json:"player_selected_answer_id"`
	QuizAnswerID           string `json:"quiz_answer_id"`
	IsCorrect              bool   `json:"is_correct"`
}

// TODO:
// - Also get score from written answers and add it with selected answers
// - Put SQL queries in their own .sql files maybe (?)

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

	scores, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerScore])
	if err != nil {
		return nil, err
	}

	var results []Result

	sql = `
	WITH player_answers AS (
		SELECT 
			player_selected_answers.player_selected_answer_id,
			player_selected_answers.user_id,
			quiz_answers.quiz_answer_id,
			quiz_answers.is_correct,
			quiz_questions.order_number,
			quiz_questions.points
		FROM player_selected_answers
		JOIN quiz_answers ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE quiz_questions.quiz_id = ($1) AND player_selected_answers.user_id = ($2)
		ORDER BY quiz_questions.order_number
	)
	SELECT 
		player_selected_answer_id, 
		quiz_answer_id, 
		is_correct
	FROM player_answers
	`

	for _, score := range scores {
		rows, err := dr.Querier.Query(ctx, sql, quizID, score.UserID)
		if err != nil {
			return nil, err
		}

		answers, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerAnswer])
		if err != nil {
			return nil, err
		}

		results = append(results, Result{score, answers})
	}

	return results, nil
}
