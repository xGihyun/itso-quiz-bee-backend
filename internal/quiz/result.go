package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Result struct {
	PlayerScore
	Answers []PlayerAnswer `json:"answers"`
}

type PlayerScore struct {
    user.GetUserResponse
	Score  int16  `json:"score"`
}

type PlayerAnswer struct {
	Answer
	PlayerAnswerID string `json:"player_answer_id"`
	QuizQuestionID string `json:"quiz_question_id"`
}

func (r *repository) GetResults(ctx context.Context, quizID string) ([]Result, error) {
	sql := `
	WITH player_written_scores AS (
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
		LEFT JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = player_written_answers.quiz_question_id
		LEFT JOIN quiz_answers 
			ON quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
		WHERE 
			quiz_questions.quiz_id = ($1)
		GROUP BY 
			player_written_answers.user_id
	)
	SELECT 
		users.user_id,
		users.created_at,
		users.username,
		users.role,
		users.name,
		COALESCE(SUM(player_score.score), 0) AS score
	FROM (
		SELECT user_id, score FROM player_written_scores
	) player_score
	RIGHT JOIN players_in_quizzes ON players_in_quizzes.user_id = player_score.user_id
	JOIN users ON users.user_id = players_in_quizzes.user_id
	GROUP BY users.user_id;
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
	WITH player_written_answers AS (
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
