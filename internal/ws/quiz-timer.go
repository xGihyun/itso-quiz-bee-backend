package ws

import (
	"time"
)

type Timer struct {
	isAuto   bool
	isPaused bool
	ticker   *time.Ticker
	duration int // Seconds
}

func (t *Timer) Start(duration int) {
	t.isPaused = false
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
	t.isPaused = false
}

func (t *Timer) Pause() {
	t.ticker.Stop()
	t.isPaused = true
}

func (t *Timer) Resume() {
	t.ticker = time.NewTicker(1 * time.Second)
	t.isPaused = false
}

// Easiest way to store timer's state
var timerState = Timer{}
