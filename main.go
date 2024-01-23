package main

import (
	"database/sql"
	"habitus/db"
	"habitus/handlers"
	"habitus/middleware"
	"habitus/models"
	"habitus/services"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stackus/dotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var habits []models.Habit
var dailys []models.Daily

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading env", "err", err)
	}
	// dbUrl := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_AUTH_TOKEN"))

	dbUrl := "file:my_db.sqlite"

	dbase, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal("error connecting to db", "err", err)
	}

	habitStore := db.NewHabitStore(dbase)
	habitService := services.NewHabit(slog.Default(), habitStore)
	habitHandler := handlers.NewHabitHandler(slog.Default(), habitService)

	dailyStore := db.NewDailyStore(slog.Default(), dbase)
	dailyService := services.NewDaily(slog.Default(), dailyStore)
	dailyHandler := handlers.NewDailyHandler(slog.Default(), dailyService)

	userStore := db.NewUserStore(slog.Default(), dbase)
	userService := services.NewUserService(slog.Default(), userStore)
	sessionStore := db.NewSessionStore(slog.Default(), dbase)
	sessionService := services.NewSessionService(slog.Default(), sessionStore, userStore)
	userHandler := handlers.NewUserHandler(slog.Default(), userService, sessionService)

	indexHandler := handlers.NewIndexHandler(habitService, dailyService)

	sessionManager := middleware.NewSessionManager(sessionService)

	router := chi.NewRouter()
	fs := http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))
	router.Use(chiMiddleware.Logger)
	router.Handle("/assets/*", fs)

	router.Group(func(r chi.Router) {
		r.Use(sessionManager.Middleware)

		r.Get("/", indexHandler.Get)

		r.Post("/{habitId}/up", habitHandler.CountUp)
		r.Post("/{habitId}/down", habitHandler.CountDown)
		r.Put("/habit", habitHandler.Put)
		r.Get("/habitModal", habitHandler.Modal)

		r.Post("/{dailyId}/done", dailyHandler.CompleteDaily)
		r.Put("/daily", dailyHandler.Put)
		r.Get("/dailyModal", dailyHandler.Modal)
		r.Get("/welcome", userHandler.GetWelcome)
	})
	router.Get("/login", userHandler.GetLogin)
	router.Post("/login", userHandler.PostLogin)
	router.Get("/signup", userHandler.GetSignup)
	router.Post("/signup", userHandler.PostSignup)
	router.Get("/logout", userHandler.Logout)

	http.ListenAndServe(":3000", router)
}
