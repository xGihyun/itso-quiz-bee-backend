# Change necessary details
DB_STRING=user=gihyun dbname=itso-quiz-bee password=password host=0.0.0.0 sslmode=disable

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_STRING)" GOOSE_MIGRATION_DIR="./internal/database/migrations/" goose up

create:
ifeq ($(name),)
	$(error `name` is not set. Usage: `make create name="migration name"`)
endif
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_STRING)" GOOSE_MIGRATION_DIR="./internal/database/migrations/" goose create "$(name)" sql
