package lobby

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, data NewLobbyRequest) (NewLobbyResponse, error)
	Join(ctx context.Context, data JoinRequest) error
	// CreateQuestion(ctx context.Context, question NewQuestion, quizID string) error
	// GetResults(ctx context.Context, quizID string) ([]Result, error)
	// CreateSelectedAnswer(ctx context.Context, data NewSelectedAnswer) error
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
