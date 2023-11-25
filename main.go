package main

import (
	"habitus/db"
	"habitus/handlers"
	"habitus/models"
	"habitus/services"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var habits []models.Habit
var dailys []models.Daily

func main() {
	habits = []models.Habit{
		{Id: 0, Name: "Water", HasDown: true},
		{Id: 1, Name: "Walk"},
		{Id: 2, Name: "Read"},
		{Id: 3, Name: "Chores"},
	}

	dailys = []models.Daily{
		{Id: 0, Name: "No Snooze", Due: true},
		{Id: 1, Name: "Read", Due: true},
	}

	habitStore := db.NewHabitStore(&habits)
	habitService := services.NewHabit(slog.Default(), habitStore)
	habitHandler := handlers.NewHabitHandler(slog.Default(), habitService)

	dailyStore := db.NewDailyStore(&dailys)
	dailyService := services.NewDaily(slog.Default(), dailyStore)
	dailyHandler := handlers.NewDailyHandler(slog.Default(), dailyService)

	indexHandler := handlers.NewIndexHandler(habitService, dailyService)

	r := chi.NewRouter()
	fs := http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))
	r.Use(middleware.Logger)
	r.Handle("/assets/*", fs)
	r.Get("/", indexHandler.Get)

	r.Post("/{habitId}/up", habitHandler.CountUp)
	r.Post("/{habitId}/down", habitHandler.CountDown)
	r.Put("/habit", habitHandler.Put)
	r.Get("/habitModal", habitHandler.Modal)

	r.Post("/{dailyId}/done", dailyHandler.CompleteDaily)
	r.Put("/daily", dailyHandler.Put)
	r.Get("/dailyModal", dailyHandler.Modal)

	http.ListenAndServe(":3000", r)
}
