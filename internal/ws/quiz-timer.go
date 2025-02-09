package ws

import (
	"context"
	"time"

	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

type Timer struct {
	isAuto   bool
	isPaused bool
	ticker   *time.Ticker
	cancel   context.CancelFunc
	duration int
}

func (c *Client) resumeQuestionTimer() {
	if c.timer.cancel != nil {
		c.timer.cancel()
	}

	timerCtx, cancel := context.WithCancel(context.Background())
	c.timer.cancel = cancel

	go func() {
		// TODO: Make sure the current question persist properly since it resets
		// when players refresh
		c.startQuestionTimer(timerCtx, c.question)
		c.timer.isPaused = false
	}()
}

func (c *Client) startQuestionTimer(ctx context.Context, question quiz.UpdatePlayersQuestionRequest) {
	if !c.timer.isPaused {
		c.timer.duration = *question.Duration
	}

	c.timer.ticker = time.NewTicker(1 * time.Second)
	defer c.timer.ticker.Stop()

	for {
		select {
		case <-c.timer.ticker.C:
			c.timer.duration -= 1

            // TODO:
            // We only need the duration
			data := quiz.QuestionTimer{
				Question:      question.Question,
				RemainingTime: c.timer.duration,
				IsAuto:        c.timer.isAuto,
			}

			response := Response{
				Event: TimerPass,
				Data:  data,
			}

			c.Pool.Broadcast <- response

			if c.timer.duration <= 0 {
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
