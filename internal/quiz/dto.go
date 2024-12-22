package quiz

type CreateSelectedAnswerRequest struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	UserID       string `json:"user_id"`
}

type CreateWrittenAnswerRequest struct {
	Content        string `json:"content"`
	QuizQuestionID string `json:"quiz_question_id"`
	UserID         string `json:"user_id"`
}

type GetWrittenAnswerResponse struct {
	PlayerWrittenAnswerID string `json:"player_written_answer_id"`
	Content               string `json:"content"`
}
