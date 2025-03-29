package quiz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
)

func (s *Service) Save(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	var data Quiz

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON request.",
		}
	}

	if err := s.repo.Save(ctx, data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "Failed to created quiz.",
		}
	}

	return api.Response{
		Code:    http.StatusCreated,
		Message: "Successfully created quiz.",
	}
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	quizID := r.PathValue("quizId")
	view := r.URL.Query().Get("view")
	includeAnswersParam := r.URL.Query().Get("includeAnswers")

	includeAnswers := false
	if includeAnswersParam == "true" {
		includeAnswers = true
	}

	var (
		result any
		err    error
	)

	switch view {
	case "basic":
		result, err = s.repo.GetBasicInfo(ctx, quizID)
	default:
		result, err = s.repo.Get(ctx, quizID, includeAnswers)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:   fmt.Errorf("get quiz: %w", err),
				Code:    http.StatusNotFound,
				Message: "Quiz not found.",
			}
		}

		return api.Response{
			Error:   fmt.Errorf("get quiz: %w", err),
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch quiz.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Data:    result,
		Message: "Fetched quiz.",
	}
}

func (s *Service) ListBasicInfo(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	results, err := s.repo.ListBasicInfo(ctx)
	if err != nil {
		return api.Response{
			Error:   fmt.Errorf("list quiz basic info: %w", err),
			Code:    http.StatusInternalServerError,
			Data:    results,
			Message: "Failed to fetch quizzes.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Data:    results,
		Message: "Fetched quizzes.",
	}
}

func (s *Service) GetPlayer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	quizID := r.PathValue("quizId")
	playerID := r.PathValue("playerId")

	request := GetPlayerRequest{UserID: playerID, QuizID: quizID}

	player, err := s.repo.GetPlayer(ctx, request)
	if err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch quiz player.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Data:    player,
		Message: "Fetched quiz player.",
	}
}

func (s *Service) GetPlayers(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	quizID := r.PathValue("quizId")

	results, err := s.repo.GetPlayers(ctx, quizID)
	if err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Data:    results,
			Message: "Failed to fetch quiz players.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Data:    results,
		Message: "Fetched quiz players.",
	}
}

func (s *Service) AddPlayer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	var data AddPlayerRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON request.",
		}
	}

	user, err := s.repo.AddPlayer(ctx, data)
	if err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "Failed to add player to quiz.",
		}
	}

	return api.Response{
		Code:    http.StatusCreated,
		Message: user.Name + " has joined.",
		Data:    user,
	}
}

func (s *Service) GetCurrentQuestion(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	quizID := r.PathValue("quizId")

	question, err := s.repo.GetCurrentQuestion(ctx, quizID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Response{
				Error:   err,
				Code:    http.StatusNotFound,
				Message: "Quiz current question not found.",
			}
		}

		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch quiz current question.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Message: "Fetched quiz current question.",
		Data:    question,
	}
}

func (s *Service) CreateWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	var data CreateWrittenAnswerRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON request.",
		}
	}

	if err := s.repo.CreateWrittenAnswer(ctx, data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "Failed to create answer: " + data.Content,
		}
	}

	return api.Response{
		Code:    http.StatusCreated,
		Message: "Successfully created answer: " + data.Content,
	}
}

func (s *Service) GetWrittenAnswer(w http.ResponseWriter, r *http.Request) api.Response {
	ctx := r.Context()

	quizID := r.PathValue("quizId")

	var data GetWrittenAnswerRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusBadRequest,
			Message: "Invalid JSON request.",
		}
	}

	if data.QuizID != quizID {
		return api.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "Invalid JSON request.",
		}
	}

	answer, err := s.repo.GetWrittenAnswer(ctx, data)
	if err != nil {
		return api.Response{
			Error:   err,
			Code:    http.StatusNotFound,
			Message: "Written answer not found.",
		}
	}

	return api.Response{
		Code:    http.StatusOK,
		Message: "Fetched written answer.",
		Data:    answer,
	}
}
