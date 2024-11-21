-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
ADD COLUMN start_at TIMESTAMP, 
ADD COLUMN end_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
