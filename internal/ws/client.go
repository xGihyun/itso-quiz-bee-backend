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

	querier database.Querier
}

func (c *Client) Read(ctx context.Context) error {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	quizRepo := quiz.NewRepository(c.querier)

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			return err
		}

		var request Request
		var response Response

		if err := json.Unmarshal(data, &request); err != nil {
			return err
		}

		response.Event = request.Event

		switch request.Event {
		case QuizUpdateStatus:
			var data quiz.UpdateStatusRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				return err
			}

			if err := quizRepo.UpdateStatus(ctx, data); err != nil {
				return err
			}

			if data.Status == quiz.Paused {
				QuizTimer.Pause()
			} else {
				QuizTimer.Resume()
				go c.handleQuestionTimer()
			}

			response.Data = data.Status

			log.Info().
				Str("event_type", string(request.Event)).
				Msg(fmt.Sprintf("Quiz status updated: %s", data.Status))

		case QuizUpdateQuestion:
			var data quiz.UpdatePlayersQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				return err
			}

			if err := quizRepo.UpdatePlayersQuestion(ctx, data); err != nil {
				return err
			}

			if data.Question.Duration != nil {
				QuizTimer.Start(*data.Question.Duration)
				go c.handleQuestionTimer()
			}

			response.Data = data.Question

			log.Info().Msg(fmt.Sprintf("Update to question #%d", data.OrderNumber))

		case PlayerTypeAnswer:
			var data quiz.CreateWrittenAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				return err
			}

			response.Data = data

		case PlayerSubmitAnswer:
			var data quiz.CreateWrittenAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				return err
			}

			if err := quizRepo.CreateWrittenAnswer(ctx, data); err != nil {
				return err
			}

			playerRequest := quiz.GetPlayerRequest{
				UserID: data.UserID,
				QuizID: data.QuizID,
			}

			player, err := quizRepo.GetPlayer(ctx, playerRequest)
			if err != nil {
				return err
			}

			response.Data = player

			log.Info().Msg("Submitted answer: " + data.Content)

		case PlayerJoin:
			var data quiz.AddPlayerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				return err
			}

			user, err := quizRepo.AddPlayer(ctx, data)
			if err != nil {
				return err
			}

			response.Data = user

			log.Info().Msg(fmt.Sprintf("%s has joined.", user.Name))

		case PlayerLeave:
		case Heartbeat:
			// Do nothing

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
		}

		c.Pool.Broadcast <- response
	}

	return nil
}
