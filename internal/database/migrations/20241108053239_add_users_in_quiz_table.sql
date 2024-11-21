-- +goose Up
-- +goose StatementBegin
CREATE TYPE quiz_status AS ENUM('open', 'started', 'paused', 'closed');

ALTER TABLE quizzes
DROP COLUMN status;

ALTER TABLE quizzes
ADD COLUMN status quiz_status NOT NULL;

CREATE TABLE IF NOT EXISTS users_in_quizzes (
	user_in_quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,

	user_id UUID NOT NULL,
	quiz_id UUID NOT NULL,
	quiz_question_id UUID, -- Current question

	FOREIGN KEY(user_id) REFERENCES users(user_id),
	FOREIGN KEY(quiz_id) REFERENCES quizzes(quiz_id),
	FOREIGN KEY(quiz_question_id) REFERENCES quiz_questions(quiz_question_id),
	UNIQUE(user_id, quiz_id)
);

ALTER TABLE users_in_lobbies
ADD CONSTRAINT users_in_lobbies_user_id_lobby_id_key UNIQUE(user_id, lobby_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
