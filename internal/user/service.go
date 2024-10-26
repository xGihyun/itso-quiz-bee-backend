package user

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, data UserRequest) error
	GetByID(ctx context.Context, userID string) (UserResponse, error)
}

type DatabaseRepository struct {
	Querier database.Querier
}

func NewDatabaseRepository(q database.Querier) *DatabaseRepository {
	return &DatabaseRepository{
		Querier: q,
	}
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
