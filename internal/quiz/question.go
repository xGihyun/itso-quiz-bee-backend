package quiz

import (
	"context"
	"encoding/json"
	"fmt"
)

type Question struct {
	QuizQuestionID string   `json:"quizQuestionId"`
	Content        string   `json:"content"`
	Points         int16    `json:"points"`
	OrderNumber    int16    `json:"orderNumber"`
	Duration       *int     `json:"duration"`
	Answers        []Answer `json:"answers,omitempty"`
}

type Answer struct {
	QuizAnswerID string `json:"quizAnswerId"`
	Content      string `json:"content"`
}

func (r *repository) GetCurrentQuestion(ctx context.Context, quizID string) (Question, error) {
	questionKey := fmt.Sprintf("quiz:%s:current_question", quizID)
	data, err := r.redisClient.Get(ctx, questionKey).Result()
	if err != nil {
		return Question{}, err
	}

	var question Question
	if err := json.Unmarshal([]byte(data), &question); err != nil {
		return Question{}, err
	}

	return question, nil
}

type setCurrentQuestionRequest struct {
	QuizID         string `json:"quizId"`
	QuizQuestionID string `json:"quizQuestionId"`
}

func (r *repository) setCurrentQuestion(
	ctx context.Context,
	data setCurrentQuestionRequest,
) (Question, error) {
	sql := `
	SELECT quiz_question_id, content, points, order_number, duration
	FROM quiz_questions
	WHERE quiz_question_id = ($1)
	`

	var question Question
	row := r.querier.QueryRow(ctx, sql, data.QuizQuestionID)
	if err := row.Scan(
		&question.QuizQuestionID,
		&question.Content,
		&question.Points,
		&question.OrderNumber,
		&question.Duration,
	); err != nil {
		return Question{}, err
	}

	questionKey := fmt.Sprintf("quiz:%s:current_question", data.QuizID)
	if err := r.redisClient.JSONSet(ctx, questionKey, "$", question).Err(); err != nil {
		return Question{}, err
	}

	return question, nil
}
