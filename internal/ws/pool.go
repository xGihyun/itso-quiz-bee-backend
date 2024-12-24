package ws

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Pool struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Request
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Request),
	}
}

func (p *Pool) Start() {
	log.Info().Msg("Starting WebSocket pool...")

	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true

			log.Info().Msg("User has connected.")

			// for client := range p.Clients {
			// 	client.Conn.WriteJSON(Request{Event: PlayerJoin})
			// }

		case client := <-p.Unregister:
			if _, ok := p.Clients[client]; ok {
				delete(p.Clients, client)

				for client := range p.Clients {
					fmt.Println(client)
					client.Conn.WriteJSON(Request{Event: PlayerLeave})
				}
			}

		case message := <-p.Broadcast:
			log.Info().Msg("Sending message to all clients...")

			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Error().Err(err).Send()
					return
				}
			}
		}
	}
}
