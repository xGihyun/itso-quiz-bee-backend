package quiz

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type interval struct {
	startAt time.Time
	endAt   time.Time
}

func NewInterval(duration time.Duration) interval {
	now := time.Now()
	return interval{
		startAt: now,
		endAt:   now.Add(duration),
	}
}

type timer struct {
	interval interval
	duration time.Duration
	isPaused bool
	started  chan bool
	done     chan bool
}

func NewTimer(duration time.Duration) *timer {
	return &timer{
		isPaused: false,
		duration: duration,
		started:  make(chan bool),
		done:     make(chan bool),
	}
}

func (t *timer) start() interval {
	interv := NewInterval(t.duration)
	t.interval = interv
	t.started <- true
	fmt.Printf("Timer start %v\n", t.interval.startAt)
	time.AfterFunc(t.duration, func() {
		t.done <- true
		fmt.Printf("Timer done %v\n", t.interval.endAt)
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

	// TODO: Cancel the previous timer before starting a new one
	timer.start()
}

func (tm *timerManager) handleTimer(quizID string) {
	timer, exists := tm.timers[quizID]
	if !exists {
		log.Error().Msg("timer not found")
		return
	}

	for {
		select {
		case <-timer.started:
			response := ws.Response{
				Event:  timerStart,
				Target: ws.All,
				Data:   timer.interval,
			}
			tm.hub.SendToAll(response)

		case <-timer.done:
			response := ws.Response{
				Event:  timerDone,
				Target: ws.All,
			}
			tm.hub.SendToAll(response)
			return
		}
	}
}
