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
	QuizChangeQuestion   Event = "quiz-change-question"
	QuizSubmitAnswer     Event = "quiz-submit-answer"
	QuizSelectAnswer     Event = "quiz-select-answer"
	QuizTypeAnswer       Event = "quiz-type-answer"
	QuizDisableAnswering Event = "quiz-disable-answering"

	QuizStartTimer Event = "quiz-start-timer"
	QuizTimerPass  Event = "quiz-timer-pass"

	PlayerJoin  Event = "player-join"
	PlayerLeave Event = "player-leave"

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

	querier database.Querier
}

type updateQuizRequest struct {
	quiz.BasicInfo
	QuizQuestionID *string `json:"quiz_question_id,omitempty"`
}

type updateQuestionRequest struct {
	QuizID string `json:"quiz_id"`
	quiz.Question
}

type addPlayerRequest struct {
	QuizID string `json:"quiz_id"`
}

type submitAnswerRequest struct {
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

		// NOTE: Is this necessary?
		request.UserID = c.ID

		switch request.Event {
		case QuizUpdateStatus:
			var data updateQuizRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdateBasicInfo(ctx, data.BasicInfo); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if data.Status == quiz.Started && data.QuizQuestionID != nil {
				if err := quizRepo.UpdatePlayersQuestion(ctx, quiz.UpdatePlayersQuestionRequest{
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
			var data updateQuestionRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.UpdatePlayersQuestion(ctx, quiz.UpdatePlayersQuestionRequest{
				QuizID:         data.QuizID,
				QuizQuestionID: data.QuizQuestionID,
			}); err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Info().Msg("Change question.")
			break

		case QuizSubmitAnswer:
			var data submitAnswerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			switch data.Variant {
			case quiz.Written:
				var answer quiz.CreateWrittenAnswerRequest

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

		case PlayerJoin:
			log.Info().Msg("User Join!")

			var data addPlayerRequest

			if err := json.Unmarshal(request.Data, &data); err != nil {
				log.Error().Err(err).Send()
				return
			}

			if err := quizRepo.AddPlayer(ctx, quiz.AddPlayerRequest{
				UserID: c.ID,
				QuizID: data.QuizID,
			}); err != nil {
				log.Error().Err(err).Send()
				return
			}

            // TODO: Use user repo instead
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
