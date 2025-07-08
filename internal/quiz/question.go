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

type currentQuestion struct {
	Question Question `json:"question"`
	Interval interval `json:"interval,omitzero"`
}

func (r *repository) GetCurrentQuestion(
	ctx context.Context,
	quizID string,
) (currentQuestion, error) {
	questionKey := fmt.Sprintf("quiz:%s:current_question", quizID)
	data, err := r.redisClient.JSONGet(ctx, questionKey).Result()
	if err != nil {
		return currentQuestion{}, err
	}

	var question Question
	if err := json.Unmarshal([]byte(data), &question); err != nil {
		return currentQuestion{}, err
	}

	intervalKey := fmt.Sprintf("quiz:%s:current_question_interval", quizID)
	data, err = r.redisClient.JSONGet(ctx, intervalKey).Result()
	if err != nil {
		return currentQuestion{}, err
	}

	var interval interval
	if err := json.Unmarshal([]byte(data), &interval); err != nil {
		return currentQuestion{}, err
	}

	curQuestion := currentQuestion{
		Question: question,
		Interval: interval,
	}

	return curQuestion, nil
}

type setCurrentQuestionRequest struct {
	QuizID         string `json:"quizId"`
	QuizQuestionID string `json:"quizQuestionId"`
}

func (r *repository) setCurrentQuestion(
	ctx context.Context,
	data setCurrentQuestionRequest,
) (currentQuestion, error) {
	sql := `
	SELECT quiz_question_id, content, points, order_number, 
		extract(epoch FROM duration)::int as duration
	FROM quiz_questions
	WHERE quiz_question_id = ($1)
	`

	var question currentQuestion
	row := r.querier.QueryRow(ctx, sql, data.QuizQuestionID)
	if err := row.Scan(
		&question.Question.QuizQuestionID,
		&question.Question.Content,
		&question.Question.Points,
		&question.Question.OrderNumber,
		&question.Question.Duration,
	); err != nil {
		return currentQuestion{}, err
	}

	questionKey := fmt.Sprintf("quiz:%s:current_question", data.QuizID)
	if err := r.redisClient.JSONSet(ctx, questionKey, "$", question).Err(); err != nil {
		return currentQuestion{}, err
	}

	return question, nil
}
