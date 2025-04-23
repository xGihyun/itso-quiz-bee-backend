package quiz

import (
	"context"
	"encoding/json"

	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type webSocketServer struct {
	repo         Repository
	timerManager *timerManager
}

func NewWebSocketServer(repo Repository, wsHub *ws.Hub) *webSocketServer {
	return &webSocketServer{
		repo:         repo,
		timerManager: NewTimerManager(wsHub),
	}
}

const (
	updateStatus    ws.Event = "quiz:update-status"
	updateQuestion  ws.Event = "quiz:update-question"
	showLeaderboard ws.Event = "quiz:show-leaderboard"

	timerStart ws.Event = "quiz:timer-start"
	timerDone  ws.Event = "quiz:timer-done"

	playerJoin         ws.Event = "quiz:player-join"
	playerLeave        ws.Event = "quiz:player-leave"
	playerTypeAnswer   ws.Event = "quiz:player-type-answer"
	playerSubmitAnswer ws.Event = "quiz:player-submit-answer"
)

func (s *webSocketServer) Handle(ctx context.Context, request ws.Request) (ws.Response, error) {
	switch request.Event {
	case updateStatus:
		var data UpdateStatusRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		if err := s.repo.UpdateStatus(ctx, data); err != nil {
			return ws.Response{}, err
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   data.Status,
		}, nil

	case updateQuestion:
		var data setCurrentQuestionRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		question, err := s.repo.setCurrentQuestion(ctx, data)
		if err != nil {
			return ws.Response{}, err
		}

		if question.Duration != nil {
			go s.timerManager.handleTimer(data.QuizID)
			s.timerManager.startTimer(data.QuizID, 5)
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   question,
		}, nil

	case showLeaderboard:
		var data bool
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}
		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   data,
		}, nil

	case playerTypeAnswer:
		var data CreateWrittenAnswerRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.Admins,
			Data:   data,
		}, nil

	case playerSubmitAnswer:
		var data CreateWrittenAnswerRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		if err := s.repo.CreateWrittenAnswer(ctx, data); err != nil {
			return ws.Response{}, err
		}

		playerRequest := GetPlayerRequest{
			UserID: data.UserID,
			QuizID: data.QuizID,
		}

		player, err := s.repo.GetPlayer(ctx, playerRequest)
		if err != nil {
			return ws.Response{}, err
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.Admins,
			Data:   player,
		}, nil

	case playerJoin:
		var data AddPlayerRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		user, err := s.repo.AddPlayer(ctx, data)
		if err != nil {
			return ws.Response{}, err
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   user,
		}, nil
	}

	return ws.Response{}, nil
}
