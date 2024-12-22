package quiz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) Create(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data Quiz

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

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	results, err := s.repo.GetAll(ctx)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{Data: results, Status: api.Success, StatusCode: http.StatusOK, Message: "Fetched all quizzes."}
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
			Data:       results,
		}
	}

	return api.Response{Data: results, Status: api.Success, StatusCode: http.StatusOK, Message: "Fetched quiz results."}
}

func (qs *Service) CreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data NewSelectedAnswer

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
				Status:     api.Fail,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
				Status:     api.Error,
			}
		}
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	data.UserID = cookie.Value

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

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
				Status:     api.Fail,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
				Status:     api.Error,
			}
		}
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	data.UserID = cookie.Value

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

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
				Status:     api.Fail,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
				Status:     api.Error,
			}
		}
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	data.UserID = cookie.Value

	if err := qs.repo.Join(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success, Message: "Joined quiz."}
}

func (s *Service) GetCurrentQuestion(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	question, err := s.repo.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
				Status:     api.Fail,
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusOK, Status: api.Success, Message: "Fetched current question.", Data: question}
}

func (s *Service) UpdateByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data BasicInfo

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := s.repo.UpdateByID(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{
		Status:     api.Success,
		StatusCode: http.StatusOK,
		Message:    "Updated quiz info.",
	}
}

func (s *Service) UpdateStatusByID(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data UpdateStatusRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Status:     api.Fail,
		}
	}

	if err := s.repo.UpdateStatusByID(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{
		Status:     api.Success,
		StatusCode: http.StatusOK,
		Message:    "Updated quiz status.",
	}
}

func (s *Service) GetAllUsers(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	users, err := s.repo.GetAllUsers(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{Data: users, Status: api.Success, StatusCode: http.StatusOK, Message: "Fetched all users in quiz."}
}

func (s *Service) GetWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
				Status:     api.Fail,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
				Status:     api.Error,
			}
		}
	}

	quizID := r.PathValue("quiz_id")

	answer, err := s.repo.GetWrittenAnswer(ctx, quizID, cookie.Value)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
			Status:     api.Fail,
			Message:    "Written answer not found.",
		}
	}

	return api.Response{
		Status:     api.Success,
		StatusCode: http.StatusOK,
		Message:    "Fetched written answer.",
		Data:       answer,
	}
}
