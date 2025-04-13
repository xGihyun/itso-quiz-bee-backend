package quiz

import (
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type timer struct {
	startAt  time.Time
	endAt    time.Time
	duration int // seconds
	isPaused bool
	done     chan bool
}

func NewTimer(duration int) *timer {
	return &timer{
		startAt:  time.Now(),
		endAt:    time.Now().Add(time.Second * time.Duration(duration)),
		isPaused: false,
		duration: duration,
		done:     make(chan bool),
	}
}

func (t *timer) start() {
	duration := time.Second * time.Duration(t.duration)
	afterTimer := time.AfterFunc(duration, func() {
		t.done <- true
	})
	defer afterTimer.Stop()
}

type timerManager struct {
	timers map[string]*timer
	hub    *ws.Hub
}

func NewTimerManager(hub *ws.Hub) *timerManager {
	return &timerManager{
		timers: make(map[string]*timer),
		hub:    hub,
	}
}

func (tm *timerManager) startTimer(quizID string, duration int) {
	timer := NewTimer(duration)
	tm.timers[quizID] = timer

	timer.start()
}

type timerPassResponse struct {
	QuizID   string `json:"quizId"`
	Duration int    `json:"duration"`
}

func (tm *timerManager) handleTimer(quizID string) {
	timer, exists := tm.timers[quizID]
	if !exists {
		return
	}

	for {
		select {
		case <-timer.done:
			response := ws.Response{
				Event:  timerDone,
				Target: ws.All,
				Data:   quizID,
			}
			tm.hub.SendToAll(response)
			return
		}
	}
}
