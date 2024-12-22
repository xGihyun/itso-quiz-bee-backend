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

type GetCurrentQuestionRequest struct {
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}

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
	JOIN users_in_quizzes ON users_in_quizzes.quiz_question_id = quiz_questions.quiz_question_id
	WHERE users_in_quizzes.quiz_id = ($1)
	`

	row := r.querier.QueryRow(ctx, sql, quizID)

	var question Question

	if err := row.Scan(&question); err != nil {
		return Question{}, err
	}

	// NOTE: Just to make sure the answers don't get leaked xD
	question.Answers = []Answer{}

	return question, nil
}
