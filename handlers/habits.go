package handlers

import (
	"habitus/components"
	"habitus/middleware"
	"habitus/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HabitHandler struct {
	Log          *slog.Logger
	HabitService HabitService
}

type HabitService interface {
	CountUp(habitId int) services.HabitAndLog
	CountDown(habitId int) services.HabitAndLog
	AddHabit(userId int, habitName string, hasUp, hasDown bool) services.HabitAndLog
	GetHabits(userId int) []services.HabitAndLog
}

func NewHabitHandler(log *slog.Logger, habitService HabitService) *HabitHandler {
	return &HabitHandler{Log: log, HabitService: habitService}
}

func (h *HabitHandler) CountUp(w http.ResponseWriter, r *http.Request) {
	habitIndex, _ := strconv.Atoi(chi.URLParam(r, "habitId"))
	habit := h.HabitService.CountUp(habitIndex)
	components.Habit(habit).Render(r.Context(), w)
}

func (h *HabitHandler) CountDown(w http.ResponseWriter, r *http.Request) {
	habitIndex, _ := strconv.Atoi(chi.URLParam(r, "habitId"))
	habit := h.HabitService.CountDown(habitIndex)
	components.Habit(habit).Render(r.Context(), w)
}

func (h *HabitHandler) Put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formHabitName := r.Form.Get("habitName")
	formHasUp := r.Form.Get("hasUp")
	formHasDown := r.Form.Get("hasDown")

	hasUp := formHasUp == "on"
	hasDown := formHasDown == "on"
	user := middleware.UserFromContext(r.Context())
	habit := h.HabitService.AddHabit(int(user.ID), formHabitName, hasUp, hasDown)
	components.Habit(habit).Render(r.Context(), w)
}

func (h *HabitHandler) Modal(w http.ResponseWriter, r *http.Request) {
	components.HabitModal().Render(r.Context(), w)
}
