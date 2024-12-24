package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

// NOTE: Everything here is very messy

type Event string

const (
	QuizUpdateStatus     Event = "quiz-update-status"
	QuizUpdateQuestion   Event = "quiz-update-question"
	QuizSubmitAnswer     Event = "quiz-submit-answer"
	QuizTypeAnswer       Event = "quiz-type-answer"
	QuizDisableAnswering Event = "quiz-disable-answering"

	QuizStartTimer Event = "quiz-start-timer"
	QuizTimerPass  Event = "quiz-timer-pass"

	PlayerJoin  Event = "player-join"
	PlayerLeave Event = "player-leave"

	Heartbeat Event = "heartbeat"
)

type Request struct {
	Event    Event           `json:"event"`
	Data     json.RawMessage `json:"data"`
	Response any             `json:"response"`
}

type Client struct {
	Pool *Pool
	Conn *websocket.Conn
	ID   string

	querier database.Querier
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	ctx := context.Background()

	quizRepo := quiz.NewRepository(c.querier)

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

		switch request.Event {
		case QuizUpdateStatus:
			var data quiz.LiveUpdateStatusRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			question, err := quizRepo.LiveUpdateStatus(ctx, data)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			request.Response = question

			log.Info().Msg(fmt.Sprintf("Quiz status updated: %s", data.Status))
			break

		case QuizUpdateQuestion:
			var data quiz.LiveUpdateQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.LiveUpdateQuestion(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}
			request.Response = data.Question

			log.Info().Msg(fmt.Sprintf("Update to question #%s", data.OrderNumber))
			break

		case QuizSubmitAnswer:
			var data quiz.LiveSubmitAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.LiveSubmitAnswer(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Info().Msg("Submitted answer: " + data.Content)
			break

		case QuizTypeAnswer:
			break

		case PlayerJoin:
			var data quiz.AddPlayerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			user, err := quizRepo.LiveAddPlayer(ctx, quiz.AddPlayerRequest{
				UserID: c.ID,
				QuizID: data.QuizID,
			})
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			request.Response = user

			log.Info().Msg(fmt.Sprintf("%s has joined.", user.Name))
			break

		case Heartbeat:
			log.Info().Msg("Heartbeat!")
			break

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
		}

		log.Info().Msg(fmt.Sprintf("Received: %s\n", request))

		c.Pool.Broadcast <- request
	}
}
