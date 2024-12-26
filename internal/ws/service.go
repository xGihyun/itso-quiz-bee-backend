package ws

import (
	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	pool    *Pool
	querier database.Querier
}

func NewService(querier database.Querier, pool *Pool) *Service {
	return &Service{
		querier: querier,
		pool:    pool,
	}
}
