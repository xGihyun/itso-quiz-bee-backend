package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetByID(ctx context.Context, quizID string) (QuizResponse, error)
	GetAll(ctx context.Context) ([]BasicInfo, error)
	UpdateByID(ctx context.Context, data BasicInfo) error
	UpdateStatusByID(ctx context.Context, data UpdateStatusRequest) error
	Create(ctx context.Context, data NewQuizRequest) error
	CreateQuestion(ctx context.Context, question NewQuestion, quizID string, orderNumber int) error
	GetResults(ctx context.Context, quizID string) ([]Result, error)
	CreateSelectedAnswer(ctx context.Context, data NewSelectedAnswer) error
	CreateWrittenAnswer(ctx context.Context, data NewWrittenAnswerRequest) error
	Join(ctx context.Context, data JoinRequest) error
	GetCurrentQuestion(ctx context.Context, quizID string) (Question, error)
	GetUser(ctx context.Context, userID string) (User, error)
	GetAllUsers(ctx context.Context, quizID string) ([]User, error)
	GetWrittenAnswer(ctx context.Context, quizID string, userID string) (GetWrittenAnswerResponse, error)
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
