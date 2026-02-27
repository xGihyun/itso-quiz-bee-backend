package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Repository interface {
	SignIn(ctx context.Context, data signInRequest) (signInResponse, error)
	Create(ctx context.Context, data createUserRequest) error
	GetByID(ctx context.Context, userID string) (UserResponse, error)
	GetAll(ctx context.Context) ([]UserResponse, error)

	generateSessionToken() (string, error)
	createSession(ctx context.Context, token, userID string) (session, error)
	ValidateSessionToken(ctx context.Context, token string) (sessionValidationResponse, error)
	invalidateSession(ctx context.Context, token, userID string) error
}

type repository struct {
	querier     database.Querier
	redisClient *redis.Client
}

func NewRepository(q database.Querier, redisClient *redis.Client) Repository {
	return &repository{
		querier:     q,
		redisClient: redisClient,
	}
}

func (r *repository) Create(ctx context.Context, data createUserRequest) error {
	passwordHash, err := hashPassword(data.Password)
	if err != nil {
		return err
	}

	sql := `
	INSERT INTO users (username, password, role, name)
	VALUES ($1, $2, $3, $4)
	`

	if _, err := r.querier.Exec(ctx, sql, data.Username, passwordHash, data.Role, data.Name); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, userID string) (UserResponse, error) {
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

	var user UserResponse

	if err := row.Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.Username,
		&user.Role,
		&user.Name,
		&user.AvatarURL,
	); err != nil {
		return UserResponse{}, err
	}

	return user, nil
}

func (r *repository) GetAll(ctx context.Context) ([]UserResponse, error) {
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
		return []UserResponse{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserResponse])
	if err != nil {
		return []UserResponse{}, err
	}

	return users, nil
}

var errInvalidPassword = errors.New("invalid password")

func (r *repository) SignIn(ctx context.Context, data signInRequest) (signInResponse, error) {
	query := `SELECT password FROM users WHERE username = ($1)`

	var hashedPassword string

	row := r.querier.QueryRow(ctx, query, data.Username)
	if err := row.Scan(&hashedPassword); err != nil {
		return signInResponse{}, err
	}

	isMatch := checkPasswordHash(data.Password, hashedPassword)
	if !isMatch {
		return signInResponse{}, errInvalidPassword
	}

	query = `
	SELECT user_id, created_at, username, name, role, avatar_url
    FROM users
	WHERE username = ($1)
    `

	rows, err := r.querier.Query(ctx, query, data.Username)
	if err != nil {
		return signInResponse{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserResponse])
	if err != nil {
		return signInResponse{}, err
	}

	token, err := r.generateSessionToken()
	if err != nil {
		return signInResponse{}, err
	}

	_, err = r.createSession(ctx, token, user.UserID)
	if err != nil {
		return signInResponse{}, err
	}

	return signInResponse{
		User:  user,
		Token: token,
	}, nil
}
