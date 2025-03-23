package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/xGihyun/itso-quiz-bee/internal/user"
)

type Request struct {
	Event Event           `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Response struct {
	Event  Event          `json:"event"`
	Data   any            `json:"data"`
	Target DelivaryTarget `json:"target"`
}

type DelivaryTarget int

const (
	All DelivaryTarget = iota
	Admins
	SenderAndAdmins
)

type client struct {
	pool *Pool
	conn *websocket.Conn
	id   string
	role user.Role

	handlers map[string]EventHandler
}

func (c *client) Read(ctx context.Context) error {
	defer func() {
		c.pool.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		var request Request

		if err := json.Unmarshal(data, &request); err != nil {
			return err
		}

		var handler EventHandler
		for prefix, h := range c.handlers {
			if strings.HasPrefix(string(request.Event), prefix) {
				handler = h
				break
			}
		}

		if handler == nil {
			log.Warn().Msg(fmt.Sprintf("no handler found for event: %s", request.Event))
			continue
		}

		response, err := handler.Handle(ctx, request)
		if err != nil {
			return err
		}

		c.pool.Broadcast <- response
	}
}

// TODO: Delete
// switch request.Event {
// case QuizUpdateStatus:
// 	var data quiz.UpdateStatusRequest
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	if err := c.quizRepo.UpdateStatus(ctx, data); err != nil {
// 		return err
// 	}
//
// 	if data.Status == quiz.Paused {
// 		QuizTimer.Pause()
// 	} else {
// 		QuizTimer.Resume()
// 		go c.handleQuestionTimer()
// 	}
//
// 	response.Data = data.Status
//
// case QuizUpdateQuestion:
// 	var data quiz.UpdatePlayersQuestionRequest
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	if err := c.quizRepo.UpdatePlayersQuestion(ctx, data); err != nil {
// 		return err
// 	}
//
// 	if data.Question.Duration != nil {
// 		QuizTimer.Start(*data.Question.Duration)
// 		go c.handleQuestionTimer()
// 	}
//
// 	response.Data = data.Question
//
// case PlayerTypeAnswer:
// 	var data quiz.CreateWrittenAnswerRequest
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	response.Data = data
//
// case PlayerSubmitAnswer:
// 	var data quiz.CreateWrittenAnswerRequest
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	if err := c.quizRepo.CreateWrittenAnswer(ctx, data); err != nil {
// 		return err
// 	}
//
// 	playerRequest := quiz.GetPlayerRequest{
// 		UserID: data.UserID,
// 		QuizID: data.QuizID,
// 	}
//
// 	player, err := c.quizRepo.GetPlayer(ctx, playerRequest)
// 	if err != nil {
// 		return err
// 	}
//
// 	response.Data = player
//
// case PlayerJoin:
// 	var data quiz.AddPlayerRequest
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	user, err := c.quizRepo.AddPlayer(ctx, data)
// 	if err != nil {
// 		return err
// 	}
//
// 	response.Data = user

// case QuizShowLeaderboard:
// 	var data bool
//
// 	if err := json.Unmarshal(request.Data, &data); err != nil {
// 		return err
// 	}
//
// 	response.Data = data
//
// case PlayerLeave:
// case Heartbeat:
// Do nothing
//
// default:
// 	log.Warn().Msg(fmt.Sprintf("Unknown request event: %s", request.Event))
// }
