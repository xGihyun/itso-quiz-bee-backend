-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
ADD COLUMN is_active BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
