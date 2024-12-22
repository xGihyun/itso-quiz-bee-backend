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
	PlayerAnswerID string  `json:"player_answer_id"`
	QuizQuestionID string  `json:"quiz_question_id"`
	QuizAnswerID   *string `json:"quiz_answer_id"`
	Content        string  `json:"content"`
	IsCorrect      bool    `json:"is_correct"`
}

func (r *repository) GetResults(ctx context.Context, quizID string) ([]Result, error) {
	sql := `
	WITH player_selected_scores AS (
		SELECT 
			SUM(
				CASE
					WHEN quiz_answers.is_correct IS TRUE 
					THEN quiz_questions.points
					ELSE 0
				END 
			) AS score,
			player_selected_answers.user_id
		FROM player_selected_answers
		JOIN quiz_answers 
			ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE 
			quiz_questions.quiz_id = ($1)
		GROUP BY 
			player_selected_answers.user_id
	),
	player_written_scores AS (
		SELECT 
			SUM(
				CASE
					WHEN LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
					THEN quiz_questions.points
					ELSE 0
				END 
			) AS score,
			player_written_answers.user_id
		FROM player_written_answers
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = player_written_answers.quiz_question_id
		LEFT JOIN quiz_answers 
			ON quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
		WHERE 
			quiz_questions.quiz_id = ($1)
		GROUP BY 
			player_written_answers.user_id
	)
	SELECT 
		user_id,
		SUM(score) AS score
	FROM (
		SELECT user_id, score FROM player_selected_scores
		UNION ALL
		SELECT user_id, score FROM player_written_scores
	) combined_scores
	GROUP BY user_id;
	`

	rows, err := r.querier.Query(ctx, sql, quizID)
	if err != nil {
		return []Result{}, err
	}

	scores, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerScore])
	if err != nil {
		return []Result{}, err
	}

	var results []Result

	sql = `
	WITH player_selected_answers AS (
		SELECT 
			player_selected_answers.player_selected_answer_id AS player_answer_id,
			quiz_answers.quiz_answer_id,
			quiz_questions.quiz_question_id,
			quiz_questions.order_number,
			quiz_questions.points,
			quiz_answers.content,
			quiz_answers.is_correct
		FROM player_selected_answers
		JOIN quiz_answers 
			ON quiz_answers.quiz_answer_id = player_selected_answers.quiz_answer_id
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = quiz_answers.quiz_question_id
		WHERE 
			quiz_questions.quiz_id = ($1) 
			AND player_selected_answers.user_id = ($2)
	),
	player_written_answers AS (
		SELECT
			player_written_answers.player_written_answer_id AS player_answer_id,
			quiz_answers.quiz_answer_id,
			quiz_questions.quiz_question_id,
			quiz_questions.order_number,
			quiz_questions.points,
			player_written_answers.content,
			CASE 
				WHEN LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
				THEN TRUE
				ELSE FALSE
			END AS is_correct
		FROM player_written_answers
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = player_written_answers.quiz_question_id
		LEFT JOIN quiz_answers 
			ON quiz_answers.quiz_question_id = player_written_answers.quiz_question_id
			   AND LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
		WHERE 
			quiz_questions.quiz_id = ($1)
			AND player_written_answers.user_id = ($2)
	)
	SELECT 
		player_answer_id, 
		quiz_answer_id,
		content,
		is_correct,
		quiz_question_id
	FROM (
		SELECT 
			player_answer_id, 
			quiz_answer_id, 
			content, 
			order_number, 
			is_correct,
			quiz_question_id
		FROM player_selected_answers
		UNION ALL
		SELECT 
			player_answer_id, 
			quiz_answer_id, 
			content, 
			order_number, 
			is_correct,
			quiz_question_id
		FROM player_written_answers
	) player_answers
	ORDER BY order_number
	`

	for _, score := range scores {
		rows, err := r.querier.Query(ctx, sql, quizID, score.UserID)
		if err != nil {
			return []Result{}, err
		}

		answers, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerAnswer])
		if err != nil {
			return []Result{}, err
		}

		results = append(results, Result{score, answers})
	}

	if results == nil {
		return []Result{}, nil
	}

	return results, nil
}
