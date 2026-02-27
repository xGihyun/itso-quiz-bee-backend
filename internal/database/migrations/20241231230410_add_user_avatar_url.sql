-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN avatar_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
