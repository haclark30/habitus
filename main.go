package main

import (
	"database/sql"
	_ "embed"
	"habitus/db_sqlc"
	"habitus/handlers"
	"habitus/middleware"
	"habitus/models"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stackus/dotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var (
	habits []models.Habit
	dailys []models.Daily
)

//go:embed schema.sql
var ddl string

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading env", "err", err)
	}
	// dbUrl := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_AUTH_TOKEN"))

	dbUrl := "file:my_db2.sqlite"

	dbase, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal("error connecting to db", "err", err)
	}

	queries := db_sqlc.New(dbase)

	habitHandler := handlers.NewHabitHandler(slog.Default(), queries)
	dailyHandler := handlers.NewDailyHandler(slog.Default(), queries)
	userHandler := handlers.NewUserHandler(slog.Default(), queries)
	indexHandler := handlers.NewIndexHandler(slog.Default(), queries)
	sessionManager := middleware.NewSessionManager(slog.Default(), queries)

	router := chi.NewRouter()
	fs := http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))
	router.Use(chiMiddleware.Logger)
	router.Handle("/assets/*", fs)

	router.Group(func(r chi.Router) {
		r.Use(sessionManager.Middleware)

		r.Get("/", indexHandler.Get)

		r.Post("/{habitId}/{habitLogId}/up", habitHandler.CountUp)
		r.Post("/{habitId}/{habitLogId}/down", habitHandler.CountDown)
		r.Put("/habit", habitHandler.Put)
		r.Get("/habitModal", habitHandler.Modal)
		r.Get("/habit/{habitId}/edit", habitHandler.Edit)

		r.Post("/{dailyId}/{dailyLogId}/done", dailyHandler.CompleteDaily)
		r.Put("/daily", dailyHandler.Put)
		r.Get("/dailyModal", dailyHandler.Modal)
	})
	router.Get("/login", userHandler.GetLogin)
	router.Post("/login", userHandler.PostLogin)
	router.Get("/signup", userHandler.GetSignup)
	router.Post("/signup", userHandler.PostSignup)
	router.Get("/logout", userHandler.Logout)

	http.ListenAndServe(":3000", router)
}
