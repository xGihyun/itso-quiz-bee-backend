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
	user user.UserResponse
	send chan Response

	handlers map[string]EventHandler
}

func (c *client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}

		var request Request
		if err := json.Unmarshal(data, &request); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		s := strings.Split(string(request.Event), ":")
		key := s[0]

		handler, ok := c.handlers[key]
		if !ok {
			log.Warn().Msg("handler not found for: " + key)
			continue
		}

		response, err := handler.Handle(ctx, request)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		switch response.Target {
		case All:
			c.hub.SendToAll(response)
		case Admins:
			c.hub.SendToRole(user.Admin, response)
		}
	}
}

func (c *client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Error().Msg("hub closed channel")
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Error().Err(fmt.Errorf("websocket write json: %w", err)).Send()
			}
		}
	}
}
