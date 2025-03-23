package ws

import (
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
	ctx := r.Context()
	conn, err := upgrade(w, r)
	defer conn.Close()

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	userId := r.URL.Query().Get("user_id")
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	client := &client{
		conn:     conn,
		pool:     s.pool,
		id:       uuid.NewString(),
		// quizRepo: s.quizRepo,
		role:     user.Role,
        handlers: s.handlers,
	}

	s.pool.Register <- client

	if err := client.Read(ctx); err != nil {
		log.Error().Err(err).Send()
	}
}
