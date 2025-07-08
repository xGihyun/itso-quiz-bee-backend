package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (s *Service) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := upgrade(w, r)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	defer conn.Close()

	token := r.URL.Query().Get("token")

	result, err := s.userRepo.ValidateSessionToken(ctx, token)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	client := NewClient(conn, s.hub, result.User, s.handlers)

	s.hub.register <- client

	go client.writePump()
	client.readPump(ctx)
}
