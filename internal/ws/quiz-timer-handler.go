package ws

import (
	"context"
	"fmt"

	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

func (c *Client) resumeQuestionTimer(ctx context.Context, quizID string) {
	quizRepo := quiz.NewRepository(c.querier)

	question, err := quizRepo.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		return
	}

	// Restores the saved state incase the client's state resets (i.e. browser refresh)
	if c.timer.ticker == nil {
		c.timer = timerState
	}

	go func() {
		c.timer.Resume()
		c.handleQuestionTimer(quizID, question)
	}()
}

func (c *Client) handleQuestionTimer(quizID string, question quiz.Question) {
	for {
		select {
		case <-c.timer.ticker.C:
			c.timer.duration -= 1
			timerState = c.timer

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
