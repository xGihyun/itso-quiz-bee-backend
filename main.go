package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/xGihyun/itso-quiz-bee/internal/api"
	"github.com/xGihyun/itso-quiz-bee/internal/auth"
	"github.com/xGihyun/itso-quiz-bee/internal/lobby"
	"github.com/xGihyun/itso-quiz-bee/internal/middleware"
	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
	"github.com/xGihyun/itso-quiz-bee/internal/ws"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Env struct {
	auth auth.Service

	user  user.Service
	quiz  quiz.Service
	ws    ws.Service
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env file.")
	}

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal().Msg("DATABASE_URL not found.")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database.")
	}

	defer pool.Close()

	wsPool := ws.NewPool()
	go wsPool.Start()

	env := &Env{
		auth:  *auth.NewService(pool),
		user:  *user.NewService(user.NewRepository(pool)),
		quiz:  *quiz.NewService(quiz.NewDatabaseRepository(pool)),
		ws:    *ws.NewService(*ws.NewDatabaseRepository(pool), wsPool),
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /ws", env.ws.HandleConnection)
	router.HandleFunc("GET /", health)

	// router.Handle("GET /api/session", api.HTTPHandler(env.auth.GetCurrentUser))
	router.Handle("POST /api/login", api.HTTPHandler(env.auth.Login))
	router.Handle("POST /api/register", api.HTTPHandler(env.auth.Register))

	router.Handle("GET /api/users/{user_id}", api.HTTPHandler(env.user.GetByID))
	router.Handle("GET /api/users", api.HTTPHandler(env.user.GetAll))
	// router.HandleFunc("POST /users", env.user.Create)

	router.Handle("POST /api/quizzes", api.HTTPHandler(env.quiz.Create))
	router.Handle("GET /api/quizzes", api.HTTPHandler(env.quiz.GetAll))
	router.Handle("GET /api/quizzes/{quiz_id}", api.HTTPHandler(env.quiz.GetByID))
	router.Handle("POST /api/quizzes/{quiz_id}", api.HTTPHandler(env.quiz.Create))
	router.Handle("PATCH /api/quizzes/{quiz_id}/status", api.HTTPHandler(env.quiz.UpdateStatusByID))
	router.Handle("POST /api/quizzes/{quiz_id}/join", api.HTTPHandler(env.quiz.Join))
	router.Handle("POST /api/quizzes/{quiz_id}/selected-answers", api.HTTPHandler(env.quiz.CreateSelectedAnswer))
	router.Handle("POST /api/quizzes/{quiz_id}/written-answers", api.HTTPHandler(env.quiz.CreateWrittenAnswer))
	router.Handle("GET /api/quizzes/{quiz_id}/results", api.HTTPHandler(env.quiz.GetResults))
	router.Handle("GET /api/quizzes/{quiz_id}/questions/current", api.HTTPHandler(env.quiz.GetCurrentQuestion))
	router.Handle("GET /api/quizzes/{quiz_id}/users", api.HTTPHandler(env.quiz.GetAllUsers))

	router.Handle("GET /api/quizzes/{quiz_id}/users/answers", api.HTTPHandler(env.quiz.GetWrittenAnswer))

	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal().Msg("HOST env not found.")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal().Msg("PORT env not found.")
	}

	frontendPort, ok := os.LookupEnv("FRONTEND_PORT")
	if !ok {
		log.Fatal().Msg("FRONTEND_PORT env not found.")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://" + host + ":" + frontendPort},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := http.Server{
		Addr:    host + ":" + port,
		Handler: corsHandler.Handler(middleware.RequestLogger(router)),
	}

	log.Info().Msg(fmt.Sprintf("Starting server on port: %s", port))

	server.ListenAndServe()
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Hello, World!")
}
