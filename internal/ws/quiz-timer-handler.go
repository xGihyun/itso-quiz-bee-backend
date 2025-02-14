package ws

import (
	"context"

	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

func (c *Client) resumeQuestionTimer(ctx context.Context, quizID string) {
	quizRepo := quiz.NewRepository(c.querier)

	question, err := quizRepo.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		return
	}

	QuizTimer.Resume()

	go func() {
		c.handleQuestionTimer(ctx, quizID, question)
	}()
}

func (c *Client) handleQuestionTimer(ctx context.Context, quizID string, question quiz.Question) {
	for {
		select {
		case <-QuizTimer.ticker.C:
			QuizTimer.duration -= 1

			response := Response{
				Event: TimerPass,
				Data:  QuizTimer.duration,
			}

			c.Pool.Broadcast <- response

			if QuizTimer.duration > 0 {
				continue
			}

			if QuizTimer.IsAuto {
				c.handleNextQuestion(ctx, quizID, question.OrderNumber)
			}
			return
		}
	}
}
