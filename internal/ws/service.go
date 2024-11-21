package ws

import (
	// "context"

	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type Service struct {
	repo DatabaseRepository
	pool *Pool
}

type Repository interface {
	// HandleConnection()
}

type DatabaseRepository struct {
	Querier database.Querier
}

func NewDatabaseRepository(q database.Querier) *DatabaseRepository {
	return &DatabaseRepository{
		Querier: q,
	}
}

func NewService(repo DatabaseRepository, pool *Pool) *Service {
	return &Service{
		repo: repo,
		pool: pool,
	}
}
