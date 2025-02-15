package quiz

import (
	"time"
)

type Timer struct {
	Ticker   *time.Ticker
	Duration int // seconds
	IsPaused bool
}

func (t *Timer) Start(duration int) {
	t.Duration = duration

	if t.Ticker != nil {
		t.Ticker.Reset(1 * time.Second)
		return
	}

	t.Ticker = time.NewTicker(1 * time.Second)
}

func (t *Timer) Stop() {
	if t.Ticker != nil {
		t.Ticker.Stop()
	}
	t.IsPaused = false
}

func (t *Timer) Pause() {
	if t.Ticker != nil {
		t.Ticker.Stop()
	}
	t.IsPaused = true
}

func (t *Timer) Resume() {
	if !t.IsPaused {
		return
	}

	t.IsPaused = false

	if t.Ticker != nil {
		t.Ticker.Reset(1 * time.Second)
		return
	}

	t.Ticker = time.NewTicker(1 * time.Second)
}
