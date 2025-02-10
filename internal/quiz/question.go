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
	Duration       *int            `json:"duration"`
	Answers        []Answer        `json:"answers"`
}

type Answer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	Content      string `json:"content"`
	IsCorrect    bool   `json:"is_correct"`
}

// NOTE:
// This assumes that all players are on the same question.
// This is used to persist the current question during an ongoing quiz in case
// the user refreshes the page.
func (r *repository) GetCurrentQuestion(ctx context.Context, quizID string) (Question, error) {
	sql := `
	SELECT 
		jsonb_build_object(
			'quiz_question_id', quiz_questions.quiz_question_id,
			'content', quiz_questions.content,
			'variant', quiz_questions.variant,
			'points', quiz_questions.points,
			'order_number', quiz_questions.order_number,
			'duration', EXTRACT(epoch FROM quiz_questions.duration)::INT,
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
	JOIN players_in_quizzes ON players_in_quizzes.quiz_question_id = quiz_questions.quiz_question_id
	WHERE players_in_quizzes.quiz_id = ($1)
	`

	row := r.querier.QueryRow(ctx, sql, quizID)

	var question Question

	if err := row.Scan(&question); err != nil {
		return Question{}, err
	}

	return question, nil
}

type GetNextQuestionRequest struct {
	QuizID      string `json:"quiz_id"`
	OrderNumber int16  `json:"order_number"`
}

func (r *repository) GetNextQuestion(ctx context.Context, data GetNextQuestionRequest) (Question, error) {
	sql := `
	SELECT 
		jsonb_build_object(
			'quiz_question_id', quiz_questions.quiz_question_id,
			'content', quiz_questions.content,
			'variant', quiz_questions.variant,
			'points', quiz_questions.points,
			'order_number', quiz_questions.order_number,
			'duration', EXTRACT(epoch FROM quiz_questions.duration)::INT,
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
	WHERE quiz_id = ($1) AND order_number = ($2)
	`

	row := r.querier.QueryRow(ctx, sql, data.QuizID, data.OrderNumber+1)

	var question Question

	if err := row.Scan(&question); err != nil {
		return Question{}, err
	}

	return question, nil
}
