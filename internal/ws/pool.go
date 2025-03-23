package ws

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Pool struct {
	Clients    map[*client]bool
	Register   chan *client
	Unregister chan *client
	Broadcast  chan Response
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*client]bool),
		Register:   make(chan *client),
		Unregister: make(chan *client),
		Broadcast:  make(chan Response),
	}
}

func (p *Pool) Start() {
	log.Info().Msg("Starting WebSocket pool...")

	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true

			log.Info().Msg("User has connected.")
			log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(p.Clients)))

		case client := <-p.Unregister:
			if _, ok := p.Clients[client]; ok {
				delete(p.Clients, client)

				for client := range p.Clients {
					client.conn.WriteJSON(Request{Event: "client:leave"})
				}

				log.Info().Msg("User has disconnected.")
				log.Info().Msg(fmt.Sprintf("Size of pool: %d", len(p.Clients)))
			}

		case message := <-p.Broadcast:
			switch message.Target {
			case All:
				for client := range p.Clients {
					if err := client.conn.WriteJSON(message); err != nil {
						log.Error().Err(err).Send()
					}
				}
			case Admins:
				for client := range p.Clients {
					// NOTE: Potential import cycle
					if client.role == user.Admin {
						if err := client.conn.WriteJSON(message); err != nil {
							log.Error().Err(err).Send()
						}
					}
				}
			}
		}
	}
}
