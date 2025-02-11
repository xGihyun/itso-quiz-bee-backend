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

	if question.Duration == nil {
		log.Warn().Msg("Question has no duration.")
		return
	}

	go func() {
		c.timer.Start(*question.Duration)
		c.handleQuestionTimer(question.QuizID, question.Question)
	}()
}

func (c *Client) handleNextQuestion(ctx context.Context, quizID string, orderNumber int16) {
	quizRepo := quiz.NewRepository(c.querier)

	request := quiz.GetNextQuestionRequest{
		QuizID:      quizID,
		OrderNumber: orderNumber,
	}

	question, err := quizRepo.GetNextQuestion(ctx, request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msg("No more questions.")
		}
        c.timer.Stop()
		return
	}

	data := quiz.UpdatePlayersQuestionRequest{
		Question: question,
		QuizID:   quizID,
	}

	c.updateQuestion(ctx, data)

	response := Response{
		Event: QuizUpdateQuestion,
		Data:  question,
	}

	c.Pool.Broadcast <- response

	log.Info().Msg(fmt.Sprintf("Update to question #%d", data.OrderNumber))
}
