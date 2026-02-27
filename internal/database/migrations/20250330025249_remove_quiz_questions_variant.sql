-- +goose Up
-- +goose StatementBegin
ALTER TABLE quiz_questions
DROP COLUMN variant;

DROP TYPE quiz_question_variant;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE quiz_question_variant AS ENUM('multiple-choice', 'boolean', 'written');

ALTER TABLE quiz_questions
ADD COLUMN variant quiz_question_variant NOT NULL;
-- +goose StatementEnd
