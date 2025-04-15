package ws

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type EventHandler interface {
	Handle(ctx context.Context, request Request) (Response, error)
}

type Service struct {
	hub      *Hub
	handlers map[string]EventHandler
	userRepo user.Repository
}

func NewService(
	hub *Hub,
	handlers map[string]EventHandler,
	userRepo user.Repository,
) *Service {
	return &Service{
		hub:      hub,
		handlers: handlers,
		userRepo: userRepo,
	}
}
