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
	userRepo user.Repository
	handlers map[string]EventHandler
}

func NewService(hub *Hub, userRepo user.Repository, handlers map[string]EventHandler) *Service {
	return &Service{
		hub:      hub,
		userRepo: userRepo,
		handlers: handlers,
	}
}
