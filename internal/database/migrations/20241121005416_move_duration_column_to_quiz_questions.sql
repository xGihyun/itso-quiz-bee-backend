-- +goose Up
-- +goose StatementBegin
ALTER TABLE quizzes
DROP COLUMN duration;

ALTER TABLE quiz_questions
ADD COLUMN duration INTERVAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
