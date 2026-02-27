-- +goose Up
-- +goose StatementBegin

-- NOTE: Change from `timestamp` to `timestamptz` as per Postgres' recommendation.
ALTER TABLE users
ADD COLUMN first_name TEXT NOT NULL,
ADD COLUMN middle_name TEXT,
ADD COLUMN last_name TEXT NOT NULL,
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE quizzes
DROP COLUMN is_active,
DROP COLUMN lobby_id,
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE player_selected_answers
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE player_written_answers
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE quiz_answers
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE quiz_questions
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

ALTER TABLE users_in_quizzes
RENAME TO players_in_quizzes;

ALTER TABLE players_in_quizzes
RENAME COLUMN user_in_quiz_id TO player_in_quiz_id;

ALTER TABLE players_in_quizzes
ALTER created_at TYPE timestamptz,
ALTER updated_at TYPE timestamptz;

DROP TABLE IF EXISTS user_details;
DROP TABLE IF EXISTS users_in_lobbies;
DROP TABLE IF EXISTS lobby_codes;
DROP TABLE IF EXISTS lobbies;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
