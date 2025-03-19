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
	"github.com/xGihyun/itso-quiz-bee/internal/middleware"
	"github.com/xGihyun/itso-quiz-bee/internal/quiz"
	"github.com/xGihyun/itso-quiz-bee/internal/user"
	"github.com/xGihyun/itso-quiz-bee/internal/ws"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	auth auth.Service

	user user.Service
	quiz quiz.Service
	ws   ws.Service
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load env file.")
	}

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal().Msg("DATABASE_URL not found.")
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal().Msg("HOST not found.")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal().Msg("PORT not found.")
	}

	frontendPort, ok := os.LookupEnv("FRONTEND_PORT")
	if !ok {
		log.Fatal().Msg("FRONTEND_PORT not found.")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database.")
	}
	defer pool.Close()

	wsPool := ws.NewPool()
	go wsPool.Start()

	app := &app{
		auth: *auth.NewService(pool),
		user: *user.NewService(user.NewRepository(pool)),
		quiz: *quiz.NewService(quiz.NewRepository(pool)),
		ws:   *ws.NewService(pool, wsPool),
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /ws", app.ws.HandleConnection)
	router.HandleFunc("GET /", health)

	// router.Handle("GET /api/session", api.HTTPHandler(app.auth.GetCurrentUser))
	router.Handle("POST /api/login", api.HTTPHandler(app.auth.Login))
	router.Handle("POST /api/register", api.HTTPHandler(app.auth.Register))

	router.Handle("GET /api/users", api.HTTPHandler(app.user.GetAll))
	router.Handle("GET /api/users/{user_id}", api.HTTPHandler(app.user.GetByID))
	// router.HandleFunc("POST /users", app.user.Create)

	router.Handle("GET /api/quizzes", api.HTTPHandler(app.quiz.GetMany))
	router.Handle("POST /api/quizzes", api.HTTPHandler(app.quiz.Create))
	router.Handle("GET /api/quizzes/{quiz_id}", api.HTTPHandler(app.quiz.GetByID))
	router.Handle("GET /api/quizzes/{quiz_id}/players", api.HTTPHandler(app.quiz.GetPlayers))
	router.Handle(
		"GET /api/quizzes/{quiz_id}/players/{player_id}",
		api.HTTPHandler(app.quiz.GetPlayer),
	)

	// TODO:
	// Hide the `answers[]` from players
	router.Handle(
		"GET /api/quizzes/{quiz_id}/current-question",
		api.HTTPHandler(app.quiz.GetCurrentQuestion),
	)

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
