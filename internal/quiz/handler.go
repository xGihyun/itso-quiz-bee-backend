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
			Message:    "Invalid JSON request.",
		}
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to created quiz.",
		}
	}

	return api.Response{
		StatusCode: http.StatusCreated,
		Message:    "Successfully created quiz.",
	}
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

func (s *Service) GetMany(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	results, err := s.repo.GetMany(ctx)
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

func (s *Service) AddPlayer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data AddPlayerRequest

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
		}
	}

	// NOTE: Is this necessary?
	data.UserID = cookie.Value

	user, err := s.repo.AddPlayer(ctx, data)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{
		StatusCode: http.StatusCreated,
		Message:    user.Name + " has joined.",
		Data:       user,
	}
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
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return api.Response{
		StatusCode: http.StatusOK,
		Message:    "Fetched current question.",
		Data:       question,
	}
}

func (s *Service) GetPlayers(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	users, err := s.repo.GetPlayers(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch players in quiz.",
		}
	}

	return api.Response{
		Data:       users,
		StatusCode: http.StatusOK,
		Message:    "Fetched all players in quiz.",
	}
}

func (s *Service) CreateSelectedAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data CreateSelectedAnswerRequest

	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Cookie not found",
				StatusCode: http.StatusBadRequest,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Server cookie error.",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	data.UserID = cookie.Value

	if err := s.repo.CreateSelectedAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create selected answer.",
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Message: "Submitted answer."}
}

func (s *Service) CreateWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data CreateWrittenAnswerRequest

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

	if err := s.repo.CreateWrittenAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Status:     api.Error,
		}
	}

	return api.Response{StatusCode: http.StatusCreated, Status: api.Success}
}

func (s *Service) GetWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	// TODO: Use `user_id`
	cookie, err := r.Cookie("session")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return api.Response{
				Error:      err,
				Message:    "Session not found",
				StatusCode: http.StatusBadRequest,
			}
		default:
			return api.Response{
				Error:      err,
				Message:    "Failed to fetch session.",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	quizID := r.PathValue("quiz_id")

	answer, err := s.repo.GetWrittenAnswer(ctx, quizID, cookie.Value)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusNotFound,
			Message:    "Written answer not found.",
		}
	}

	return api.Response{
		StatusCode: http.StatusOK,
		Message:    "Fetched written answer.",
		Data:       answer,
	}
}
