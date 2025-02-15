package quiz

import (
	"time"
)

type BasicInfo struct {
	QuizID      string    `json:"quiz_id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
}

type Quiz struct {
	BasicInfo
	Questions []Question `json:"questions"`
}
