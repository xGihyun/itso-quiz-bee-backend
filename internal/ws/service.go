package ws

import (
	"context"
)

type EventHandler interface {
	Handle(ctx context.Context, request Request) (Response, error)
}

type Service struct {
	hub       *Hub
	handlers  map[string]EventHandler
	jwtSecret string
}

func NewService(
	hub *Hub,
	handlers map[string]EventHandler,
	jwtSecret string,
) *Service {
	return &Service{
		hub:       hub,
		handlers:  handlers,
		jwtSecret: jwtSecret,
	}
}
