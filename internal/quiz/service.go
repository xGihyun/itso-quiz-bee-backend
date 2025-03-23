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

func NewSocketService(repo Repository, wsPool *ws.Pool) *SocketService {
	return &SocketService{
		repo:         repo,
		timerManager: NewTimerManager(wsPool),
	}
}

const (
	updateStatus   ws.Event = "quiz:update-status"
	updateQuestion ws.Event = "quiz:update-question"

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
		var data UpdatePlayersQuestionRequest
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return ws.Response{}, err
		}

		if err := s.repo.UpdatePlayersQuestion(ctx, data); err != nil {
			return ws.Response{}, err
		}

		// s.timerManager.StopTimer(data.QuizID)
		if data.Duration != nil {
			s.timerManager.StartTimer(ctx, data.QuizID, *data.Duration)
		}

		return ws.Response{
			Event:  request.Event,
			Target: ws.All,
			Data:   data.Question,
		}, nil
	}

	return ws.Response{}, nil
}
