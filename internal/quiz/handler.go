package quiz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) HandleCreate(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewQuiz

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	s.repo.Create(ctx, data)

	return api.Response{StatusCode: http.StatusCreated}
}

func (s *Service) HandleGetResults(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	// TODO:
	// - Also get score from written answers and add it with selected answers
	// - Put SQL queries in their own .sql files maybe (?)

	results, err := s.repo.GetResults(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := api.WriteJSON(w, results); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{}
}

func (qs *Service) HandleCreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewSelectedAnswer

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := qs.repo.CreateSelectedAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{StatusCode: http.StatusCreated}
}
