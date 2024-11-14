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

type QuizStartRequest struct {
	quiz.UpdateStatusRequest
	QuizQuestionID string `json:"quiz_question_id"`
}

type QuizNextQuestionRequest struct {
	QuizID         string `json:"quiz_id"`
	QuizQuestionID string `json:"quiz_question_id"`
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
		quizRepo := quiz.NewDatabaseRepository(c.repo.Querier)

		switch request.Event {
		case QuizStart:
			var data QuizStartRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdateStatusByID(ctx, quiz.UpdateStatusRequest{QuizID: data.QuizID, Status: data.Status}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			// TODO:
			// Get all users in quiz
			// Set the quiz_question to the quiz' first question
			// Redirect everyone to the question

			quizRepo.UpdateCurrentQuestion(ctx, quiz.UpdateCurrentQuestionRequest{
				QuizID:         data.QuizID,
				QuizQuestionID: data.QuizQuestionID,
			})

			log.Info().Msg("Quiz has started.")
			break

		case QuizNextQuestion:
		case QuizPreviousQuestion:
			var data QuizNextQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			quizRepo.UpdateCurrentQuestion(ctx, quiz.UpdateCurrentQuestionRequest{
				QuizID:         data.QuizID,
				QuizQuestionID: data.QuizQuestionID,
			})

			log.Info().Msg("Next question.")
			break

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %v", request.Event))
		}

		log.Info().Msg(fmt.Sprintf("Received: %v\n", request))

		c.Pool.Broadcast <- request

	}
}
