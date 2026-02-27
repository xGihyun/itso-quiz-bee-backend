-- +goose Up
-- +goose StatementBegin
ALTER TABLE quiz_answers
DROP COLUMN is_correct;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE quiz_answers
ADD COLUMN is_correct boolean NOT NULL;
-- +goose StatementEnd
