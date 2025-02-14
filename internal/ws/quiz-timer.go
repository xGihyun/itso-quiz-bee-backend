package ws

import (
	"time"
)

type Timer struct {
	IsAuto   bool
	IsPaused bool
	ticker   *time.Ticker
	duration int // seconds
}

// TODO: Move this to somewhere more appropriate
type UpdateTimerModeRequest struct {
	QuizID      string `json:"quiz_id"`
	IsTimerAuto bool   `json:"is_timer_auto"`
}

func (t *Timer) Start(duration int) {
	t.duration = duration

	if t.ticker != nil {
		t.ticker.Reset(1 * time.Second)
		return
	}

	t.ticker = time.NewTicker(1 * time.Second)
}

func (t *Timer) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.IsPaused = false
}

func (t *Timer) Pause() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.IsPaused = true
}

func (t *Timer) Resume() {
	if !t.IsPaused {
		return
	}

	t.IsPaused = false

	if t.ticker != nil {
		t.ticker.Reset(1 * time.Second)
		return
	}

	t.ticker = time.NewTicker(1 * time.Second)
}

// WARN: This will only work if there's only one quiz in progress.
// Otherwise, the existing timer state would be overwritten.
// TODO: Sync with the timer state stored on the database.
var QuizTimer Timer 

// NOTE:
// It's much better to have a separate timer per quiz
// var quizTimer map[string]Timer
