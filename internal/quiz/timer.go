package quiz

import (
	"context"
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type Timer struct {
	Ticker   *time.Ticker
	Duration int // seconds
	IsPaused bool
}

func NewTimer(duration int) *Timer {
	return &Timer{
		Ticker:   time.NewTicker(time.Second),
		Duration: duration,
		IsPaused: false,
	}
}

func (t *Timer) Stop() {
	t.Ticker.Stop()
	t.IsPaused = false
}

func (t *Timer) Pause() {
	t.Ticker.Stop()
	t.IsPaused = true
}

func (t *Timer) Start() {
	if !t.IsPaused {
		return
	}

	t.Ticker = time.NewTicker(time.Second)
	t.IsPaused = false
}

func (t *Timer) Resume() {
	t.Ticker.Reset(time.Second)
	t.IsPaused = false
}

type TimerManager struct {
	timers map[string]*Timer
	hub  *ws.Hub
}

func NewTimerManager(hub *ws.Hub) *TimerManager {
	return &TimerManager{
		timers: make(map[string]*Timer),
		hub:  hub,
	}
}

func (tm *TimerManager) StartTimer(ctx context.Context, quizID string, duration int) {
	timer, exists := tm.timers[quizID]
	if !exists {
		timer = NewTimer(duration)
		tm.timers[quizID] = timer
	}

	timer.Start()
	go tm.handleTimer(ctx, quizID)
}

func (tm *TimerManager) StopTimer(quizID string) {
	if timer, exists := tm.timers[quizID]; exists {
		timer.Stop()
		delete(tm.timers, quizID)
	}
}

func (tm *TimerManager) PauseTimer(quizID string) {
	if timer, exists := tm.timers[quizID]; exists {
		timer.Pause()
	}
}

func (tm *TimerManager) ResumeTimer(ctx context.Context, quizID string) {
	if timer, exists := tm.timers[quizID]; exists {
		timer.Resume()
		go tm.handleTimer(ctx, quizID)
	}
}

type timerPassResponse struct {
	QuizID   string `json:"quizId"`
	Duration int    `json:"duration"`
}

func (tm *TimerManager) handleTimer(ctx context.Context, quizID string) {
	timer, exists := tm.timers[quizID]
	if !exists {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.Ticker.C:
			timer.Duration -= 1

			tpResponse := timerPassResponse{
				QuizID:   quizID,
				Duration: timer.Duration,
			}
			response := ws.Response{
				Event:  timerPass,
				Target: ws.All,
				Data:   tpResponse,
			}

			tm.hub.Broadcast <- response

			if timer.Duration <= 0 {
				doneResponse := ws.Response{
					Event:  timerDone,
					Target: ws.All,
					Data:   quizID,
				}

				tm.hub.Broadcast <- doneResponse

				tm.StopTimer(quizID)
				return
			}
		}
	}
}
