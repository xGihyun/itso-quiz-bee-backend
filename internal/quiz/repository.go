package quiz

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Repository interface {
	GetByID(ctx context.Context, quizID string) (Quiz, error)
	GetMany(ctx context.Context) ([]BasicInfo, error)
	Save(ctx context.Context, data Quiz) error

	GetCurrentQuestion(ctx context.Context, quizID string) (Question, error)
	GetNextQuestion(ctx context.Context, data GetNextQuestionRequest) (Question, error)

	GetWrittenAnswer(ctx context.Context, data GetWrittenAnswerRequest) (GetWrittenAnswerResponse, error)
	CreateSelectedAnswer(ctx context.Context, data CreateSelectedAnswerRequest) error
	CreateWrittenAnswer(ctx context.Context, data CreateWrittenAnswerRequest) error

	AddPlayer(ctx context.Context, data AddPlayerRequest) (user.UserResponse, error)
	GetPlayer(ctx context.Context, data GetPlayerRequest) (Player, error)
	GetPlayers(ctx context.Context, quizID string) ([]Player, error)

	UpdateStatus(ctx context.Context, data UpdateStatusRequest) error
	Start(ctx context.Context, quizID string) (Question, error)
	UpdatePlayersQuestion(ctx context.Context, data UpdatePlayersQuestionRequest) error
}

type repository struct {
	querier database.Querier
}

func NewRepository(q database.Querier) Repository {
	return &repository{
		querier: q,
	}
}
