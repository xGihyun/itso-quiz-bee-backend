-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN middle_name,
DROP COLUMN last_name,
ADD COLUMN name TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
