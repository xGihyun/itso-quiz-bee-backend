package ws

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Hub struct {
	clients       map[*client]bool
	clientsByRole map[user.Role]map[*client]bool
	register      chan *client
	unregister    chan *client
	Broadcast     chan Response
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*client]bool),
		clientsByRole: make(map[user.Role]map[*client]bool),
		register:      make(chan *client),
		unregister:    make(chan *client),
		Broadcast:     make(chan Response),
	}
}

func (h *Hub) Start() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientsByRole[client.role][client] = true

			log.Info().Msg("User has connected.")
			log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(h.clients)))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clientsByRole[client.role], client)
				delete(h.clients, client)

				log.Info().Msg("User has disconnected.")
				log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(h.clients)))
			}
		}
	}
}

func (h *Hub) sendToRole(role user.Role, response Response) {
	clients := h.clientsByRole[role]

	for client := range clients {
		client.send <- response
	}
}
