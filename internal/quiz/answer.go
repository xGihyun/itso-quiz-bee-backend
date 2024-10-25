package quiz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

type NewAnswer struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

func (d Dependency) CreateAnswer(ctx context.Context, answer NewAnswer, questionID string) error {
	sql := `
	INSERT INTO quiz_answers (content, is_correct, quiz_question_id)
	VALUES ($1, $2, $3)
    `

	if _, err := d.DB.Exec(ctx, sql, answer.Content, answer.IsCorrect, questionID); err != nil {
		return err
	}

	return nil
}

type NewSelectedAnswer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	UserID       string `json:"user_id"`
}

func (d Dependency) CreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewSelectedAnswer

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
	INSERT INTO player_selected_answers (quiz_answer_id, user_id)
	VALUES ($1, $2)
    `

	if _, err := d.DB.Exec(ctx, sql, data.QuizAnswerID, data.UserID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
