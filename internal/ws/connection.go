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

	// keepAlive(conn, time.Duration(25)*time.Second)
}

// func keepAlive(conn *websocket.Conn, timeout time.Duration) {
// 	lastResponse := time.Now()
//
// 	conn.SetPongHandler(func(msg string) error {
// 		lastResponse = time.Now()
// 		// log.Debug().Msg("Received pong from client!")
// 		return nil
// 	})
//
// 	go func() {
// 		for {
// 			err := conn.WriteMessage(websocket.PingMessage, []byte("Ping!"))
// 			if err != nil {
// 				log.Err(err).Msg("Failed to write ping message.")
// 				return
// 			}
//
// 			time.Sleep(timeout / 2)
// 			if time.Now().Sub(lastResponse) > timeout {
// 				log.Warn().Msg(fmt.Sprintf("No ping response, disconnecting to %s", conn.LocalAddr()))
// 				_ = conn.Close()
// 				return
// 			}
// 		}
// 	}()
// }
