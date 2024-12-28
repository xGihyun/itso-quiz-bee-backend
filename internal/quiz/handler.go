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
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:      err,
				StatusCode: http.StatusNotFound,
				Message:    "Quiz not found.",
			}
		}

		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch quiz.",
		}
	}

	return api.Response{
		StatusCode: http.StatusOK,
		Data:       result,
		Message:    "Fetched quiz.",
	}
}

func (s *Service) GetMany(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	results, err := s.repo.GetMany(ctx)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Data:       results,
			Message:    "Failed to fetch quizzes.",
		}
	}

	return api.Response{
		StatusCode: http.StatusOK,
		Data:       results,
		Message:    "Fetched all quizzes.",
	}
}

func (s *Service) GetResults(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	results, err := s.repo.GetResults(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Data:       results,
			Message:    "Failed to fetch quiz results.",
		}
	}

	return api.Response{
		StatusCode: http.StatusOK,
		Data:       results,
		Message:    "Fetched quiz results.",
	}
}

func (s *Service) AddPlayer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data AddPlayerRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	user, err := s.repo.AddPlayer(ctx, data)
	if err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to add player to quiz.",
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

func (s *Service) CreateWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	var data CreateWrittenAnswerRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	if err := s.repo.CreateWrittenAnswer(ctx, data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create answer: " + data.Content,
		}
	}

	return api.Response{
		StatusCode: http.StatusCreated,
		Message:    "Successfully created answer: " + data.Content,
	}
}

func (s *Service) GetWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := context.Background()

	quizID := r.PathValue("quiz_id")

	var data GetWrittenAnswerRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:      err,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON request.",
		}
	}

	if data.QuizID != quizID {
		return api.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Invalid JSON request.",
		}
	}

	answer, err := s.repo.GetWrittenAnswer(ctx, data)
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
