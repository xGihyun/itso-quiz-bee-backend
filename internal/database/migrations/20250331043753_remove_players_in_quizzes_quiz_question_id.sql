-- +goose Up
-- +goose StatementBegin
ALTER TABLE players_in_quizzes
DROP COLUMN quiz_question_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players_in_quizzes
ADD COLUMN quiz_question_id uuid NOT NULL;

ALTER TABLE players_in_quizzes
ADD CONSTRAINT players_in_quizzes_quiz_question_id_fk uuid NOT NULL;
-- +goose StatementEnd
