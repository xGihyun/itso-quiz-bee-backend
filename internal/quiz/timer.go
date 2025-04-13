package quiz

import (
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type interval struct {
	startAt time.Time
	endAt   time.Time
}

func NewInterval(duration time.Duration) interval {
	return interval{
		startAt: time.Now(),
		endAt:   time.Now().Add(time.Second * duration),
	}
}

type timer struct {
	interval interval
	duration time.Duration
	isPaused bool
	done     chan bool
}

func NewTimer(duration time.Duration) *timer {
	return &timer{
		isPaused: false,
		duration: duration,
		done:     make(chan bool),
	}
}

func (t *timer) start() interval {
	interv := NewInterval(t.duration)
	time.AfterFunc(t.duration, func() {
		t.done <- true
	})
	return interv
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
	dur := time.Second * time.Duration(duration)
	timer := NewTimer(dur)
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
