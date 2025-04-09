package ws

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Hub struct {
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	Broadcast  chan Response
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		Broadcast:  make(chan Response),
	}
}

func (h *Hub) Start() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			log.Info().Msg("User has connected.")
			log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(h.clients)))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

				for client := range h.clients {
					client.conn.WriteJSON(Request{Event: "client:leave"})
				}

				log.Info().Msg("User has disconnected.")
				log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(h.clients)))
			}

			// case message := <-h.Broadcast:
			// 	switch message.Target {
			// 	case All:
			// 		for client := range h.clients {
			// 			if err := client.conn.WriteJSON(message); err != nil {
			// 				log.Error().Err(err).Send()
			// 			}
			// 		}
			//
			// 	case Admins:
			// 		for client := range h.clients {
			// 			// NOTE: Potential import cycle
			// 			if client.role == user.Admin {
			// 				if err := client.conn.WriteJSON(message); err != nil {
			// 					log.Error().Err(err).Send()
			// 				}
			// 			}
			// 		}
			// 	}
		}
	}
}
