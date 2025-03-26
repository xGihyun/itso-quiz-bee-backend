package quiz

import (
	"time"
)

type BasicInfo struct {
	QuizID      string    `json:"quizId"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
}

type Quiz struct {
	BasicInfo
	Questions []Question `json:"questions"`
}
