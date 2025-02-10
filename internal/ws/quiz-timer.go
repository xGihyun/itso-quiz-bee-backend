package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

// TODO: Fix bug where current remaining time doesn't persist when player refreshes 
// browser while timer is paused

type Timer struct {
	isAuto   bool
	isPaused bool
	ticker   *time.Ticker
	duration int // Seconds
}

func (t *Timer) Start(duration int) {
	t.isPaused = false
	t.duration = duration
	t.ticker = time.NewTicker(1 * time.Second)
}

func (t *Timer) Stop() {
	t.ticker.Stop()
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

func (c *Client) resumeQuestionTimer(ctx context.Context, quizID string) {
	quizRepo := quiz.NewRepository(c.querier)

	question, err := quizRepo.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		return
	}

	go func() {
		c.timer.Resume()
		c.handleQuestionTimer(quizID, question)
	}()
}

func (c *Client) handleQuestionTimer(quizID string, question quiz.Question) {
	defer c.timer.Stop()

	for {
		select {
		case <-c.timer.ticker.C:
			c.timer.duration -= 1

			fmt.Println(fmt.Sprintf("TICK %d", c.timer.duration))

			response := Response{
				Event: TimerPass,
				Data:  c.timer.duration,
			}

			c.Pool.Broadcast <- response

			if c.timer.duration > 0 {
				continue
			}

			if c.timer.isAuto {
				ctx := context.Background()
				c.handleNextQuestion(ctx, quizID, question.OrderNumber)
			}
			return
		}
	}
}
