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

type Question struct {
	QuizQuestionID string          `json:"quiz_question_id"`
	Content        string          `json:"content"`
	Variant        QuestionVariant `json:"variant"`
	Points         int16           `json:"points"`
	OrderNumber    int16           `json:"order_number"`
	Answers        []Answer        `json:"answers"`
}

type NewQuestion struct {
	Content string          `json:"content"`
	Variant QuestionVariant `json:"variant"`
	Points  int16           `json:"points"`
	// OrderNumber int16           `json:"order_number"`
	Answers []NewAnswer `json:"answers"`
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

type GetCurrentQuestionRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

func (dr *DatabaseRepository) GetCurrentQuestion(ctx context.Context, quizID string) (Question, error) {
	sql := `
	SELECT 
		jsonb_build_object(
			'quiz_question_id', quiz_questions.quiz_question_id,
			'content', quiz_questions.content,
			'variant', quiz_questions.variant,
			'points', quiz_questions.points,
			'order_number', quiz_questions.order_number,
			'answers', (
				SELECT jsonb_agg(
					jsonb_build_object(
						'quiz_answer_id', quiz_answers.quiz_answer_id,
						'content', quiz_answers.content,
						'is_correct', quiz_answers.is_correct
					)
				)
				FROM quiz_answers
				WHERE quiz_answers.quiz_question_id = quiz_questions.quiz_question_id
			)
		) AS question
	FROM quiz_questions
	JOIN users_in_quizzes ON users_in_quizzes.quiz_question_id = quiz_questions.quiz_question_id
	WHERE users_in_quizzes.quiz_id = ($1)
	`

	row := dr.Querier.QueryRow(ctx, sql, quizID)

	var question Question

	if err := row.Scan(&question); err != nil {
		return Question{}, err
	}

	return question, nil
}
