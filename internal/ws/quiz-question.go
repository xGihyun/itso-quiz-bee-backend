package ws

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
)

func (c *Client) updateQuestion(ctx context.Context, question quiz.UpdatePlayersQuestionRequest) {
	quizRepo := quiz.NewRepository(c.querier)

	if err := quizRepo.UpdatePlayersQuestion(ctx, question); err != nil {
		log.Error().Err(err).Send()
		return
	}

	if c.timer.cancel != nil {
		c.timer.cancel()
	}

	timerCtx, cancel := context.WithCancel(context.Background())
	c.timer.cancel = cancel

    c.timer.isPaused = false
    c.question = question

	go c.startQuestionTimer(timerCtx, question)
}

func (c *Client) handleNextQuestion(ctx context.Context, question quiz.UpdatePlayersQuestionRequest) {
	quizRepo := quiz.NewRepository(c.querier)

	nextQuestion, err := quizRepo.GetNextQuestion(ctx, question)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msg("No more questions.")
		}
		return
	}

	data := quiz.UpdatePlayersQuestionRequest{
		Question: nextQuestion,
		QuizID:   question.QuizID,
	}

	c.updateQuestion(ctx, data)

	response := Response{
		Event: QuizUpdateQuestion,
		Data:  nextQuestion,
	}

	c.Pool.Broadcast <- response

	log.Info().Msg(fmt.Sprintf("Update to question #%d", data.OrderNumber))
}
