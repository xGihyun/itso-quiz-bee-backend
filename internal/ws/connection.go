package ws

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
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

	token := r.URL.Query().Get("token")

	claims, err := s.verifyToken(token)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	client := &client{
		conn:     conn,
		hub:      s.hub,
		role:     claims.Role,
		handlers: s.handlers,
	}

	s.hub.register <- client

	if err := client.Read(ctx); err != nil {
		log.Error().Err(err).Send()
	}
}

type userClaims struct {
	jwt.RegisteredClaims

	UserID string    `json:"userId"`
	Role   user.Role `json:"role"`
}

func (s *Service) verifyToken(tokenString string) (*userClaims, error) {
	claims := userClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			return s.jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
