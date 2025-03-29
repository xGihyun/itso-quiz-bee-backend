-- +goose Up
-- +goose StatementBegin
DROP TABLE player_selected_answers;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS player_selected_answers (
	player_selected_answer_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	quiz_answer_id uuid NOT NULL,
	user_id uuid NOT NULL,

    FOREIGN KEY(quiz_answer_id) REFERENCES quiz_answers(quiz_answer_id),
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);
-- +goose StatementEnd
