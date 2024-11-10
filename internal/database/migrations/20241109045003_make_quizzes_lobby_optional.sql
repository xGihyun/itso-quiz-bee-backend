-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
ALTER COLUMN lobby_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
