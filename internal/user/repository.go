package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Repository interface {
	Create(ctx context.Context, data CreateUserRequest) error
	GetByID(ctx context.Context, userID string) (GetUserResponse, error)
	GetAll(ctx context.Context) ([]GetUserResponse, error)
}

type repository struct {
	querier database.Querier

	// Insert other dependencies if needed ...
}

func NewRepository(q database.Querier) Repository {
	return &repository{
		querier: q,
	}
}

// TODO:
// - Password hashing
func (r *repository) Create(ctx context.Context, data CreateUserRequest) error {
	sql := `
	INSERT INTO users (username, password, role, name)
	VALUES ($1, $2, $3, $4)
	`

	if _, err := r.querier.Exec(ctx, sql, data.Username, data.Password, data.Role, data.Name); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, userID string) (GetUserResponse, error) {
	sql := `
    SELECT 
        user_id, 
        created_at,
        username,
        role,
        name,
        avatar_url
    FROM users WHERE user_id = ($1)
    `

	row := r.querier.QueryRow(ctx, sql, userID)

	var user GetUserResponse

	if err := row.Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.Username,
		&user.Role,
		&user.Name,
		&user.AvatarURL,
	); err != nil {
		return GetUserResponse{}, err
	}

	return user, nil
}

func (r *repository) GetAll(ctx context.Context) ([]GetUserResponse, error) {
	query := `
    SELECT 
        user_id, 
        created_at,
        username,
        role,
        name,
        avatar_url
	FROM users
	`

	rows, err := r.querier.Query(ctx, query)
	if err != nil {
		return []GetUserResponse{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[GetUserResponse])
	if err != nil {
		return []GetUserResponse{}, err
	}

	return users, nil
}
