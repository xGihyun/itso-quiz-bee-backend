package ws

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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
	conn, err := upgrade(w, r)
	defer conn.Close()

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	client := &Client{
		Conn:        conn,
		Pool:        s.pool,
		ID:          uuid.NewString(),
		querier:     s.querier,
	}

	s.pool.Register <- client

	ctx := context.Background()

	client.Read(ctx)
}
