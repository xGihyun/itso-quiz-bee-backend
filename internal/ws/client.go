package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Event string

type Request struct {
	Event Event           `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Response struct {
	Event  Event          `json:"event"`
	Data   any            `json:"data"`
	Target DelivaryTarget `json:"-"`
}

type DelivaryTarget int

const (
	All DelivaryTarget = iota
	Admins
	SenderAndAdmins
)

type client struct {
	hub  *Hub
	conn *websocket.Conn
	role user.Role

	handlers map[string]EventHandler
}

func (c *client) Read(ctx context.Context) error {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		var request Request
		if err := json.Unmarshal(data, &request); err != nil {
			return err
		}

		s := strings.Split(":", string(request.Event))
		key := s[0]

		handler, ok := c.handlers[key]
		if !ok {
			log.Warn().Msg(fmt.Sprintf("no handler found for event: %s", request.Event))
			continue
		}

		response, err := handler.Handle(ctx, request)
		if err != nil {
			return err
		}

		c.hub.Broadcast <- response
	}
}
