package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

type Event string

const (
	QuizStart            Event = "quiz-start"
	QuizPause            Event = "quiz-pause"
	QuizResume           Event = "quiz-resume"
	QuizEnd              Event = "quiz-end"
	QuizNextQuestion     Event = "quiz-next-question"
	QuizPreviousQuestion Event = "quiz-previous-question"
	UserJoin             Event = "user-join"
	UserLeave            Event = "user-leave"
)

type Request struct {
	Event Event           `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Client struct {
	Pool *Pool
	Conn *websocket.Conn
	ID   string

	repo DatabaseRepository
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	ctx := context.Background()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}

		var request Request

		if err := json.Unmarshal(data, &request); err != nil {
			log.Error().Err(err).Send()
			return
		}

		// TODO: Do stuff based on request

		switch request.Event {
		case QuizStart:
			var data quiz.UpdateStatusRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			quizRepo := quiz.NewDatabaseRepository(c.repo.Querier)

			if err := quizRepo.UpdateStatusByID(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Info().Msg("Quiz has started.")
			break

		case QuizNextQuestion:
			log.Info().Msg("Next question.")
			break
		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %v", request.Event))
		}

		log.Info().Msg(fmt.Sprintf("Received: %v\n", request))

		c.Pool.Broadcast <- request

	}
}
