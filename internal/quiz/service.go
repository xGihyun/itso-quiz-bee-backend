package quiz

import (
	"context"
	"encoding/json"

	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type SocketService struct {
	repo         Repository
	timerManager *TimerManager
}

func NewSocketService(repo Repository, wsHub *ws.Hub) *SocketService {
	return &SocketService{
		repo:         repo,
		timerManager: NewTimerManager(wsHub),
	}
}

const (
	updateStatus    ws.Event = "quiz:update-status"
	updateQuestion  ws.Event = "quiz:update-question"
	showLeaderboard ws.Event = "quiz:show-leaderboard"

	timerPass ws.Event = "quiz:timer-pass"
	timerDone ws.Event = "quiz:timer-done"

	playerJoin         ws.Event = "quiz:player-join"
	playerLeave        ws.Event = "quiz:player-leave"
	playerTypeAnswer   ws.Event = "quiz:player-type-answer"
	playerSubmitAnswer ws.Event = "quiz:player-submit-answer"
)

func (s *SocketService) Handle(ctx context.Context, request ws.Request) (ws.Response, error) {
	switch request.Event {
	case updateStatus:
		var data UpdateStatusRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		if err := s.repo.UpdateStatus(ctx, data); err != nil {
			return ws.Response{}, err
		}

		if data.Status == Paused {
			s.timerManager.PauseTimer(data.QuizID)
		} else if data.Status == Started {
			s.timerManager.ResumeTimer(ctx, data.QuizID)
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

		// s.timerManager.StopTimer(data.QuizID)
		if question.Duration != nil {
			s.timerManager.StartTimer(ctx, data.QuizID, *question.Duration)
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   data,
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
			Target: ws.All,
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
			Target: ws.All,
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
