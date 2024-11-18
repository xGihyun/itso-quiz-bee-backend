package user

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Role string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type UserResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
	Detail
}

type Detail struct {
	FirstName  string  `json:"first_name"`
	MiddleName *string `json:"middle_name"`
	LastName   string  `json:"last_name"`
}

// TODO:
// - Password hashing?
// - Not sure if this should be on auth
func (dr *DatabaseRepository) Create(ctx context.Context, data UserRequest) error {
	sql := `
	INSERT INTO users (email, password, role)
	VALUES ($1, $2, $3)
	`

	if _, err := dr.Querier.Exec(ctx, sql, data.Email, data.Password, data.Role); err != nil {
		return err
	}

	return nil
}

func (dr *DatabaseRepository) GetByID(ctx context.Context, userID string) (UserResponse, error) {
	query := "SELECT user_id, email, role FROM users WHERE user_id = ($1)"

	row := dr.Querier.QueryRow(ctx, query, userID)

	var user UserResponse

	if err := row.Scan(&user.UserID, &user.Email, &user.Role); err != nil {
		return UserResponse{}, err
	}

	return user, nil
}

func (dr *DatabaseRepository) GetAll(ctx context.Context) ([]UserResponse, error) {
	query := `
	SELECT 
		users.user_id, 
		users.email, 
		users.role,
		user_details.first_name,
		user_details.middle_name,
		user_details.last_name
	FROM users
	JOIN user_details ON user_details.user_id = users.user_id
	`

	rows, err := dr.Querier.Query(ctx, query)
	if err != nil {
		return []UserResponse{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserResponse])
	if err != nil {
		return []UserResponse{}, err
	}

	return users, nil
}
