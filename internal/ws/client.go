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

type Request struct {
	Event Event           `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Response struct {
	Event Event `json:"event"`
	Data  any   `json:"data"`
}

type Client struct {
	Pool *Pool
	Conn *websocket.Conn
	ID   string

	querier  database.Querier
	timer    Timer
}

func (c *Client) Read(ctx context.Context) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	quizRepo := quiz.NewRepository(c.querier)

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}

		var request Request
		var response Response

		if err := json.Unmarshal(data, &request); err != nil {
			log.Error().Err(err).Send()
			return
		}

		response.Event = request.Event

		switch request.Event {
		// TODO: Merge this with `QuizStart`
		case QuizUpdateStatus:
			var data quiz.UpdateStatusRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdateStatus(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			// TODO: Could create a separate event for this.
			if data.Status == quiz.Paused {
				c.timer.Pause()
			} else {
				c.resumeQuestionTimer(ctx, data.QuizID)
			}

			response.Data = data.Status

			log.Info().
				Str("event_type", string(request.Event)).
				Msg(fmt.Sprintf("Quiz status updated: %s", data.Status))

		case QuizStart:
			var quizID string

			if err := json.Unmarshal(request.Data, &quizID); err != nil {
				log.Error().Err(err).Send()
				return
			}

			question, err := quizRepo.Start(ctx, quizID)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = question

			log.Info().Msg("Quiz has started.")

		case QuizUpdateQuestion:
			var data quiz.UpdatePlayersQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			c.updateQuestion(ctx, data)
			response.Data = data.Question

			log.Info().Msg(fmt.Sprintf("Update to question #%d", data.OrderNumber))

		case PlayerTypeAnswer:
			var data quiz.CreateWrittenAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = data

		case PlayerSubmitAnswer:
			var data quiz.CreateWrittenAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.CreateWrittenAnswer(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			playerRequest := quiz.GetPlayerRequest{
				UserID: data.UserID,
				QuizID: data.QuizID,
			}

			player, err := quizRepo.GetPlayer(ctx, playerRequest)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = player

			log.Info().Msg("Submitted answer: " + data.Content)

		case PlayerJoin:
			var data quiz.AddPlayerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			user, err := quizRepo.AddPlayer(ctx, data)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = user

			log.Info().Msg(fmt.Sprintf("%s has joined.", user.Name))

		case TimerUpdateMode:
			var isAuto bool

			if err := json.Unmarshal(request.Data, &isAuto); err != nil {
				log.Error().Err(err).Send()
				return
			}

			c.timer.isAuto = isAuto

			log.Info().Msg(fmt.Sprintf("Toggled timer auto mode: %t", isAuto))

		case PlayerLeave:
		case Heartbeat:
			// Do nothing

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
		}

		c.Pool.Broadcast <- response
	}
}
