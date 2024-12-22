package quiz

import (
	"time"
)

type Status string

const (
	Open    Status = "open"
	Started Status = "started"
	Paused  Status = "paused"
	Closed  Status = "closed"
)

type BasicInfo struct {
	QuizID      string  `json:"quiz_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      Status  `json:"status"`
}

type Quiz struct {
	QuizID      string     `json:"quiz_id"`
	CreatedAt   time.Time  `json:"created_at"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      Status     `json:"status"`
	Questions   []Question `json:"questions"`
}

type Question struct {
	QuizQuestionID string          `json:"quiz_question_id"`
	Content        string          `json:"content"`
	Variant        QuestionVariant `json:"variant"`
	Points         int16           `json:"points"`
	OrderNumber    int16           `json:"order_number"`
	Duration       *time.Duration  `json:"duration"`
	Answers        []Answer        `json:"answers"`
}

type Answer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	Content      string `json:"content"`
	IsCorrect    bool   `json:"is_correct"`
}
