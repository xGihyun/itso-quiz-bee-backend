package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	keepAlive(conn, time.Duration(25)*time.Second)

	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}

		message := string(bytes[:])
		log.Info().Msg("Received: " + message)

		if err := conn.WriteMessage(messageType, bytes); err != nil {
			log.Error().Err(err).Send()
			return
		}
	}
}

func keepAlive(conn *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()

	conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		log.Debug().Msg("Received pong from client!")
		return nil
	})

	go func() {
		for {
			err := conn.WriteMessage(websocket.PingMessage, []byte("Ping!"))
			if err != nil {
				log.Err(err).Msg("Failed to write ping message.")
				return
			}

			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				log.Warn().Msg(fmt.Sprintf("No ping response, disconnecting to %s", conn.LocalAddr()))
				_ = conn.Close()
				return
			}
		}
	}()
}
