package quiz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewQuizRequest

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

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success, Message: "Quiz created."}
}

func (s *Service) GetByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	result, err := s.repo.GetByID(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{Data: result, Status: api.Success, StatusCode: http.StatusOK, Message: "Fetched quiz."}
}

func (s *Service) GetResults(w http.ResponseWriter, r *http.Request) api.Response {
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

func (qs *Service) CreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
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

func (qs *Service) CreateWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewWrittenAnswerRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := qs.repo.CreateWrittenAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success}
}

func (qs *Service) Join(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data JoinRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := qs.repo.Join(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success, Message: "Joined quiz."}
}
