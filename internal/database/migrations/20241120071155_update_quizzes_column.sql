-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
DROP COLUMN start_at,
DROP COLUMN end_at,
ADD COLUMN duration INTERVAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
