package quiz

import (
	"context"
)

type QuestionVariant string

const (
	MultipleChoice QuestionVariant = "multiple-choice"
	Boolean        QuestionVariant = "boolean"
	Written        QuestionVariant = "written"
)

type NewQuestion struct {
	Content     string          `json:"content"`
	Variant     QuestionVariant `json:"variant"`
	Points      int16           `json:"points"`
	// OrderNumber int16           `json:"order_number"`
	Answers     []NewAnswer     `json:"answers"`
}

func (dr *DatabaseRepository) CreateQuestion(ctx context.Context, question NewQuestion, quizID string, orderNumber int) error {
	sql := `
	INSERT INTO quiz_questions (content, variant, points, order_number, quiz_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING quiz_question_id
	`

	row := dr.Querier.QueryRow(ctx, sql, question.Content, question.Variant, question.Points, orderNumber, quizID)

	var questionID string

	if err := row.Scan(&questionID); err != nil {
		return err
	}

	for _, answer := range question.Answers {
		if err := dr.CreateAnswer(ctx, answer, questionID); err != nil {
			return err
		}
	}

	return nil
}
