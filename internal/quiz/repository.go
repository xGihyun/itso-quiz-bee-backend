package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Repository interface {
	GetByID(ctx context.Context, quizID string) (Quiz, error)
	GetMany(ctx context.Context) ([]BasicInfo, error)
	Create(ctx context.Context, data Quiz) error
	UpdateBasicInfo(ctx context.Context, data BasicInfo) error

	GetResults(ctx context.Context, quizID string) ([]Result, error)
	GetCurrentQuestion(ctx context.Context, quizID string) (Question, error)

	GetWrittenAnswer(ctx context.Context, quizID string, userID string) (GetWrittenAnswerResponse, error)
	CreateSelectedAnswer(ctx context.Context, data CreateSelectedAnswerRequest) error
	CreateWrittenAnswer(ctx context.Context, data CreateWrittenAnswerRequest) error

	AddPlayer(ctx context.Context, data AddPlayerRequest) error
	GetPlayers(ctx context.Context, quizID string) ([]user.GetUserResponse, error)

	LiveUpdateStatus(ctx context.Context, data LiveUpdateStatusRequest) (Question, error)
	LiveUpdateQuestion(ctx context.Context, data LiveUpdateQuestionRequest) error
	LiveSubmitAnswer(ctx context.Context, data LiveSubmitAnswerRequest) error
	LiveAddPlayer(ctx context.Context, data AddPlayerRequest) (user.GetUserResponse, error)
}

type repository struct {
	querier database.Querier
}

func NewRepository(q database.Querier) Repository {
	return &repository{
		querier: q,
	}
}
