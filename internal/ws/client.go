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
	QuizUpdateStatus     Event = "quiz-update-status"
	// QuizStart            Event = "quiz-start"
	// QuizPause            Event = "quiz-pause"
	// QuizResume           Event = "quiz-resume"
	// QuizEnd              Event = "quiz-end"
	QuizChangeQuestion   Event = "quiz-change-question"
	QuizSubmitAnswer     Event = "quiz-submit-answer"
	QuizSelectAnswer     Event = "quiz-select-answer"
	QuizTypeAnswer       Event = "quiz-type-answer"
	QuizDisableAnswering Event = "quiz-disable-answering"

	QuizStartTimer Event = "quiz-start-timer"
	QuizTimerPass  Event = "quiz-timer-pass"

	UserJoin  Event = "user-join"
	UserLeave Event = "user-leave"
	Heartbeat Event = "heartbeat"
)

type Request struct {
	Event  Event           `json:"event"`
	Data   json.RawMessage `json:"data"`
	UserID string          `json:"user_id"`
}

type Client struct {
	Pool *Pool
	Conn *websocket.Conn
	ID   string

	repo DatabaseRepository
}

type QuizUpdateStatusRequest struct {
	quiz.UpdateStatusRequest
	QuizQuestionID *string `json:"quiz_question_id,omitempty"`
}

// type QuizStartRequest struct {
// 	quiz.UpdateStatusRequest
// 	QuizQuestionID string `json:"quiz_question_id"`
// }

type QuizChangeQuestionRequest struct {
	QuizID string `json:"quiz_id"`
	quiz.Question
}

type QuizJoinRequest struct {
	QuizID string `json:"quiz_id"`
}

type QuizSubmitAnswerRequest struct {
	UserID         string               `json:"user_id"`
	QuizID         string               `json:"quiz_id"`
	QuizQuestionID string               `json:"quiz_question_id"`
	Variant        quiz.QuestionVariant `json:"variant"`
	Answer         json.RawMessage      `json:"answer"`
}

type (
	QuizSelectAnswerRequest quiz.NewSelectedAnswer
	QuizTypeAnswerRequest   quiz.NewWrittenAnswerRequest
)

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	ctx := context.Background()

	quizRepo := quiz.NewDatabaseRepository(c.repo.Querier)

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

		request.UserID = c.ID

		switch request.Event {
		case QuizUpdateStatus:
			var data QuizUpdateStatusRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdateStatusByID(ctx, quiz.UpdateStatusRequest{QuizID: data.QuizID, Status: data.Status}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if data.Status == quiz.Started && data.QuizQuestionID != nil {
				if err := quizRepo.UpdateCurrentQuestion(ctx, quiz.UpdateCurrentQuestionRequest{
					QuizID:         data.QuizID,
					QuizQuestionID: *data.QuizQuestionID,
				}); err != nil {
					log.Error().Err(err).Send()
					return
				}
			}

			log.Info().Msg(fmt.Sprintf("Quiz status updated: %s", data.Status))
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
			case quiz.MultipleChoice, quiz.Boolean:
				var answer quiz.NewSelectedAnswer

				if err := json.Unmarshal(data.Answer, &answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				answer.UserID = c.ID

				if err := quizRepo.CreateSelectedAnswer(ctx, answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				log.Info().Msg("Submitted multiple choice answer.")
				break
			case quiz.Written:
				var answer quiz.NewWrittenAnswerRequest

				if err := json.Unmarshal(data.Answer, &answer); err != nil {
					log.Error().Err(err).Send()
					return
				}

				answer.UserID = c.ID

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

		case QuizSelectAnswer:
			log.Info().Msg("Answer selected.")
			break

		case QuizTypeAnswer:
			log.Info().Msg("Answer typed.")
			break

		case UserJoin:
			log.Info().Msg("User Join!")

			var data QuizJoinRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.Join(ctx, quiz.JoinRequest{
				UserID: c.ID,
				QuizID: data.QuizID,
			}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			user, err := quizRepo.GetUser(ctx, c.ID)
			if err != nil {
				log.Error().Err(err).Send()
			}

			log.Info().Msg(fmt.Sprintf("%s", user.FirstName))

			msg, err := json.Marshal(user)
			if err != nil {
				log.Error().Err(err).Send()
			}
			request.Data = msg

			log.Info().Msg(fmt.Sprintf("%s %s has joined.", user.FirstName, user.LastName))
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
