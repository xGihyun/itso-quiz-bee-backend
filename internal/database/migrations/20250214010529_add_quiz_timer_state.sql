-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
-- NOTE: Create a new table for each quiz timer if we need more data (state)
ADD COLUMN is_timer_auto BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
