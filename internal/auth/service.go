package auth

import "github.com/xGihyun/itso-quiz-bee/internal/database"

type Service struct {
	querier database.Querier
}

func NewService(querier database.Querier) *Service {
	return &Service{
		querier: querier,
	}
}
