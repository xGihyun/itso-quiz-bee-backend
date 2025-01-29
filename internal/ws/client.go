package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	querier     database.Querier
	cancelTimer context.CancelFunc
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

			err := quizRepo.UpdateStatus(ctx, data)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = data.Status

			log.Info().Msg(fmt.Sprintf("Quiz status updated: %s", data.Status))
			break

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
			break

		case QuizUpdateQuestion:
			var data quiz.UpdatePlayersQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdatePlayersQuestion(ctx, data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if c.cancelTimer != nil {
				c.cancelTimer()
			}

			timerCtx, cancel := context.WithCancel(context.Background())
			c.cancelTimer = cancel

			go c.startQuestionTimer(timerCtx, data.Question)

			response.Data = data.Question

			log.Info().Msg(fmt.Sprintf("Update to question #%d", data.OrderNumber))
			break

		case PlayerTypeAnswer:
			var data quiz.CreateWrittenAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			response.Data = data
			break

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
			break

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
			break

		case PlayerLeave:
		case Heartbeat:
			// Do nothing
			break

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
		}

		// log.Info().Msg(fmt.Sprintf("Received: %s\n", request))

		c.Pool.Broadcast <- response
	}
}

type QuestionTimer struct {
	quiz.Question
	RemainingTime int `json:"remaining_time"`
}

func (c *Client) startQuestionTimer(ctx context.Context, question quiz.Question) {
	duration := *question.Duration

	if duration <= 0 {
		log.Warn().Msg("Invalid timer duration.")
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			duration -= 1

			data := QuestionTimer{
				Question:      question,
				RemainingTime: duration,
			}

			response := Response{
				Event: QuizTimerPass,
				Data:  data,
			}

			c.Pool.Broadcast <- response

			if duration <= 0 {
				log.Info().Msg("Time is up!")
				return
			}

		case <-ctx.Done():
			log.Info().Msg("Timer cancelled.")
			return
		}
	}
}
