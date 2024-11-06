package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, data NewQuiz) error
	CreateQuestion(ctx context.Context, question NewQuestion, quizID string) error
	GetResults(ctx context.Context, quizID string) ([]Result, error)
	CreateSelectedAnswer(ctx context.Context, data NewSelectedAnswer) error
	CreateWrittenAnswer(ctx context.Context, data NewWrittenAnswer) error
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
