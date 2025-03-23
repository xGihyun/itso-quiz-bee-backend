package ws

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type EventHandler interface {
	Handle(ctx context.Context, request Request) (Response, error)
}

type Service struct {
	pool     *Pool
	userRepo user.Repository
	handlers map[string]EventHandler
}

func NewService(pool *Pool, userRepo user.Repository, handlers map[string]EventHandler) *Service {
	return &Service{
		pool:     pool,
		userRepo: userRepo,
		handlers: handlers,
	}
}
