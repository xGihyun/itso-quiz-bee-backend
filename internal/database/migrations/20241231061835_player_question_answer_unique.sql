-- +goose Up
-- +goose StatementBegin
ALTER TABLE player_written_answers
ADD CONSTRAINT player_written_answers_quiz_question_id_user_id_unique
UNIQUE (quiz_question_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
