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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Env struct {
	user       user.Dependency
	auth       auth.Dependency
	lobby      lobby.Dependency
	quiz       quiz.Dependency
	middleware middleware.Dependency
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

	env := &Env{
		user:       user.Dependency{DB: pool},
		auth:       auth.Dependency{DB: pool},
		lobby:      lobby.Dependency{DB: pool},
		quiz:       quiz.Dependency{DB: pool},
		middleware: middleware.Dependency{Log: log.Logger},
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /", health)

	router.Handle("POST /login", api.HTTPHandler(env.auth.Login))
	router.Handle("POST /register", api.HTTPHandler(env.auth.Register))

	router.Handle("GET /users/{id}", api.HTTPHandler(env.user.GetByID))
	// router.HandleFunc("POST /users", env.user.Create)

	router.Handle("POST /lobbies", api.HTTPHandler(env.lobby.Create))
	router.Handle("POST /lobbies/join", api.HTTPHandler(env.lobby.Join))

	router.Handle("POST /quizzes", api.HTTPHandler(env.quiz.Create))
	router.Handle("POST /quizzes/answers", api.HTTPHandler(env.quiz.CreateSelectedAnswer))
	router.Handle("GET /quizzes/{quiz_id}/results", api.HTTPHandler(env.quiz.GetResults))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal().Msg("PORT not found.")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: env.middleware.RequestLogger(router),
	}

	log.Info().Msg(fmt.Sprintf("Starting server on port: %s", port))

	server.ListenAndServe()
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Hello, World!")
}
