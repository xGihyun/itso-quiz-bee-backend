package ws

import (
	"context"
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

type Timer struct {
	isAuto        bool
	ticker        *time.Ticker
	cancel        context.CancelFunc
	remainingTime int
}

func (t *Timer) Stop() {
	t.cancel()
	t.ticker.Stop()
}

func (t *Timer) Resume() {
    t.ticker.Stop()
	t.ticker = time.NewTicker(1 * time.Second)
}

func (c *Client) startQuestionTimer(ctx context.Context, question quiz.UpdatePlayersQuestionRequest) {
	c.timer.remainingTime = 5 // TODO: Change to the question's duration

	c.timer.ticker = time.NewTicker(1 * time.Second)
	defer c.timer.ticker.Stop()

	for {
		select {
		case <-c.timer.ticker.C:
			c.timer.remainingTime -= 1

			data := quiz.QuestionTimer{
				Question:      question.Question,
				RemainingTime: c.timer.remainingTime,
				IsAuto:        c.timer.isAuto,
			}

			response := Response{
				Event: TimerPass,
				Data:  data,
			}

			c.Pool.Broadcast <- response

			if c.timer.remainingTime <= 0 {
				if c.timer.isAuto {
					c.handleNextQuestion(ctx, question)
				}
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
