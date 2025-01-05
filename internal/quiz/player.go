package quiz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type AddPlayerRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (r *repository) AddPlayer(ctx context.Context, data AddPlayerRequest) (user.GetUserResponse, error) {
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

type Player struct {
	user.GetUserResponse
	Result PlayerResult `json:"result"`
}

type PlayerResult struct {
	Score   int16          `json:"score"`
	Answers []PlayerAnswer `json:"answers"`
}

type GetPlayerRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (r *repository) GetPlayer(ctx context.Context, data GetPlayerRequest) (Player, error) {
	sql := `
	WITH player_scores AS (
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
	),
	player_answers AS (
		SELECT
            jsonb_agg(
                jsonb_build_object(
                    'player_answer_id', player_written_answers.player_written_answer_id,
                    'quiz_question_id', quiz_questions.quiz_question_id,
                    'order_number', quiz_questions.order_number,
                    'points', quiz_questions.points,
                    'content', player_written_answers.content,
                    'is_correct', LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
                )
            ) AS answers,
            player_written_answers.user_id
		FROM player_written_answers
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = player_written_answers.quiz_question_id
		LEFT JOIN quiz_answers 
            ON quiz_answers.quiz_question_id = player_written_answers.quiz_question_id
            AND LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
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
        users.avatar_url,
        jsonb_build_object(
            'score', COALESCE(SUM(player_scores.score), 0),
            'answers', COALESCE(player_answers.answers, '[]'::jsonb)
        ) AS result
    FROM players_in_quizzes
    LEFT JOIN player_answers ON player_answers.user_id = players_in_quizzes.user_id
    LEFT JOIN player_scores ON player_scores.user_id = players_in_quizzes.user_id
    JOIN users ON users.user_id = players_in_quizzes.user_id
    WHERE users.user_id = ($2)
	GROUP BY users.user_id, player_answers.answers
    `

	var player Player

	row := r.querier.QueryRow(ctx, sql, data.QuizID, data.UserID)

	if err := row.Scan(
		&player.UserID,
		&player.CreatedAt,
		&player.Username,
		&player.Role,
		&player.Name,
		&player.AvatarURL,
		&player.Result,
	); err != nil {
		return Player{}, err
	}

	return player, nil
}

func (r *repository) GetPlayers(ctx context.Context, quizID string) ([]Player, error) {
	sql := `
	WITH player_scores AS (
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
	),
	player_answers AS (
		SELECT
            jsonb_agg(
                jsonb_build_object(
                    'player_answer_id', player_written_answers.player_written_answer_id,
                    'quiz_question_id', quiz_questions.quiz_question_id,
                    'order_number', quiz_questions.order_number,
                    'points', quiz_questions.points,
                    'content', player_written_answers.content,
                    'is_correct', LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
                )
            ) AS answers,
            player_written_answers.user_id
		FROM player_written_answers
		JOIN quiz_questions 
			ON quiz_questions.quiz_question_id = player_written_answers.quiz_question_id
		LEFT JOIN quiz_answers 
            ON quiz_answers.quiz_question_id = player_written_answers.quiz_question_id
            AND LOWER(TRIM(quiz_answers.content)) = LOWER(TRIM(player_written_answers.content))
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
        users.avatar_url,
        jsonb_build_object(
            'score', COALESCE(SUM(player_scores.score), 0),
            'answers', COALESCE(player_answers.answers, '[]'::jsonb)
        ) AS result
    FROM players_in_quizzes
    LEFT JOIN player_answers ON player_answers.user_id = players_in_quizzes.user_id
    LEFT JOIN player_scores ON player_scores.user_id = players_in_quizzes.user_id
    JOIN users ON users.user_id = players_in_quizzes.user_id
	GROUP BY users.user_id, player_answers.answers
    `

	rows, err := r.querier.Query(ctx, sql, quizID)
	if err != nil {
		return []Player{}, err
	}

	players, err := pgx.CollectRows(rows, pgx.RowToStructByName[Player])
	if err != nil {
		return []Player{}, err
	}

	return players, nil
}
