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
	QuizStart          Event = "quiz-start"
	QuizPause          Event = "quiz-pause"
	QuizResume         Event = "quiz-resume"
	QuizEnd            Event = "quiz-end"
	QuizChangeQuestion Event = "quiz-change-question"
	QuizSubmitAnswer   Event = "quiz-submit-answer"
	UserJoin           Event = "user-join"
	UserLeave          Event = "user-leave"
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

type QuizChangeQuestionRequest struct {
	QuizID         string `json:"quiz_id"`
	QuizQuestionID string `json:"quiz_question_id"`
}

type QuizSubmitAnswerRequest struct {
	UserID         string               `json:"user_id"`
	QuizID         string               `json:"quiz_id"`
	QuizQuestionID string               `json:"quiz_question_id"`
	Variant        quiz.QuestionVariant `json:"variant"`
	Answer         json.RawMessage      `json:"answer"`
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

			if err := quizRepo.UpdateCurrentQuestion(ctx, quiz.UpdateCurrentQuestionRequest{
				QuizID:         data.QuizID,
				QuizQuestionID: data.QuizQuestionID,
			}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Info().Msg("Quiz has started.")
			break

		case QuizChangeQuestion:
			var data QuizChangeQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdateCurrentQuestion(ctx, quiz.UpdateCurrentQuestionRequest{
				QuizID:         data.QuizID,
				QuizQuestionID: data.QuizQuestionID,
			}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Info().Msg("Change question.")
			break

		case QuizSubmitAnswer:
			var data QuizSubmitAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			switch data.Variant {
			case quiz.MultipleChoice:
			case quiz.Boolean:
				var answer quiz.NewSelectedAnswer

				answer.UserID = c.ID

				if err := json.Unmarshal(data.Answer, &answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				if err := quizRepo.CreateSelectedAnswer(ctx, answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				log.Info().Msg("Submitted multiple choice answer.")
				break
			case quiz.Written:
				var answer quiz.NewWrittenAnswerRequest

				answer.UserID = c.ID

				if err := json.Unmarshal(data.Answer, &answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				if err := quizRepo.CreateWrittenAnswer(ctx, answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				log.Info().Msg("Submitted written answer.")
				break

			default:
				log.Warn().Msg("Invalid question variant.")
			}

			break

		default:
			log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
		}

		log.Info().Msg(fmt.Sprintf("Received: %s\n", request))

		c.Pool.Broadcast <- request

	}
}