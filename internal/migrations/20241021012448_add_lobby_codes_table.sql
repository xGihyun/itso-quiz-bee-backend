-- +goose Up
-- +goose StatementBegin
ALTER TABLE lobbies 
DROP COLUMN code;

CREATE TABLE lobby_codes (
	lobby_code_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,

	code TEXT NOT NULL UNIQUE,
	lobby_id UUID NOT NULL UNIQUE,

	FOREIGN KEY(lobby_id) REFERENCES lobbies(lobby_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
