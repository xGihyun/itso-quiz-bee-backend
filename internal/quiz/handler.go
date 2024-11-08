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
			Status:     api.Fail,
		}
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success}
}

func (s *Service) HandleGetResults(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	results, err := s.repo.GetResults(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{Data: results, Status: api.Success, StatusCode: http.StatusOK}
}

func (qs *Service) HandleCreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewSelectedAnswer

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := qs.repo.CreateSelectedAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success}
}
