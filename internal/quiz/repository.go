package quiz

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Repository interface {
	Get(ctx context.Context, quizID string, includeAnswers bool) (Quiz, error)
	GetBasicInfo(ctx context.Context, quizID string) (BasicInfo, error)
	ListBasicInfo(ctx context.Context) ([]BasicInfo, error)
	Save(ctx context.Context, data Quiz) error
	UpdateStatus(ctx context.Context, data UpdateStatusRequest) error

	GetCurrentQuestion(ctx context.Context, quizID string) (Question, error)
	setCurrentQuestion(ctx context.Context, data setCurrentQuestionRequest) (Question, error)
	GetNextQuestion(ctx context.Context, data GetNextQuestionRequest) (Question, error)

	GetWrittenAnswer(
		ctx context.Context,
		data GetWrittenAnswerRequest,
	) (GetWrittenAnswerResponse, error)
	CreateSelectedAnswer(ctx context.Context, data CreateSelectedAnswerRequest) error
	CreateWrittenAnswer(ctx context.Context, data CreateWrittenAnswerRequest) error

	AddPlayer(ctx context.Context, data AddPlayerRequest) (user.UserResponse, error)
	GetPlayer(ctx context.Context, data GetPlayerRequest) (Player, error)
	GetPlayers(ctx context.Context, quizID string) ([]Player, error)
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
