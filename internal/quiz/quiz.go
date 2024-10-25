package quiz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/database"
)

type QuizService struct {
	store Dependency
	repo QuizRepo
}

type QuizRepo interface {
	CreateQuestion(ctx context.Context, question NewQuestion, quizID string) error
}

type Dependency struct {
	DB database.Querier
}

type Status string

const (
	Open    Status = "open"
	Ongoing Status = "ongoing"
	Closed  Status = "closed"
)

type NewQuiz struct {
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Status      Status        `json:"status"`
	LobbyID     string        `json:"lobby_id"`
	Questions   []NewQuestion `json:"questions"`
}

func (d Dependency) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewQuiz

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	sql := `
    INSERT INTO quizzes (name, description, status, lobby_id)
    VALUES ($1, $2, $3, $4)
    RETURNING quiz_id
    `

	// tx, err := d.DB.Begin(ctx)
	// if err != nil {
	// 	return api.Response{
	// 		Error:      err,
	// 		StatusCode: http.StatusInternalServerError,
	// 	}
	// }

	row := d.DB.QueryRow(ctx, sql, data.Name, data.Description, data.Status, data.LobbyID)

	var quizID string

	if err := row.Scan(&quizID); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	for _, question := range data.Questions {
		if err := d.CreateQuestion(ctx, question, quizID); err != nil {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	// if err := tx.Commit(ctx); err != nil {
	// 	return api.Response{
	// 		Error:      err,
	// 		StatusCode: http.StatusInternalServerError,
	// 	}
	// }

	return api.Response{StatusCode: http.StatusCreated}
}

func (qs *QuizService) CreateInDB(ctx context.Context, data NewQuiz) error {
	sql := `
    INSERT INTO quizzes (name, description, status, lobby_id)
    VALUES ($1, $2, $3, $4)
    RETURNING quiz_id
    `

	row := qs.store.DB.QueryRow(ctx, sql, data.Name, data.Description, data.Status, data.LobbyID)

	var quizID string

	if err := row.Scan(&quizID); err != nil {
		return err
	}

	for _, question := range data.Questions {
		if err := qs.repo.CreateQuestion(ctx, question, quizID); err != nil {
			return err
		}
	}

	return nil
}
