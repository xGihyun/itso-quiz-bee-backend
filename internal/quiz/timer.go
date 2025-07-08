package quiz

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/ws"
)

type interval struct {
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
}

func NewInterval(duration time.Duration) interval {
	now := time.Now()
	return interval{
		StartAt: now,
		EndAt:   now.Add(duration),
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
	fmt.Printf("Timer start %v\n", t.interval.StartAt)
	time.AfterFunc(t.duration, func() {
		t.done <- true
		fmt.Printf("Timer done %v\n", t.interval.EndAt)
	})
	return interv
}

type timerManager struct {
	timers      map[string]*timer
	hub         *ws.Hub
	redisClient *redis.Client
}

func NewTimerManager(hub *ws.Hub, redisClient *redis.Client) *timerManager {
	return &timerManager{
		timers:      make(map[string]*timer),
		hub:         hub,
		redisClient: redisClient,
	}
}

func (tm *timerManager) startTimer(quizID string, duration int) {
	dur := time.Second * time.Duration(duration)
	timer := NewTimer(dur)
	tm.timers[quizID] = timer

	// TODO: Cancel the previous timer before starting a new one
	timer.start()
}

func (tm *timerManager) handleTimer(ctx context.Context, quizID string) {
	timer, exists := tm.timers[quizID]
	if !exists {
		log.Error().Msg("timer not found")
		return
	}

	for {
		select {
		// TODO: Store `timer.interval` on Redis to persist on client refresh
		case <-timer.started:
			intervalKey := fmt.Sprintf("quiz:%s:current_question_interval", quizID)
			if err := tm.redisClient.JSONSet(ctx, intervalKey, "$", timer.interval).Err(); err != nil {
				log.Error().Err(err).Send()
				continue
			}
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
