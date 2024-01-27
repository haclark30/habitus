package handlers

import (
	"habitus/components"
	"habitus/models"
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
	CountUp(habitId int) models.Habit
	CountDown(habitId int) models.Habit
	AddHabit(userId int, habitName string, hasDown bool) models.Habit
	GetHabits(userId int) []models.Habit
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
	formHasDown := r.Form.Get("hasDown")
	hasDown := formHasDown == "on"
	user := r.Context().Value("user").(models.User)
	habit := h.HabitService.AddHabit(user.Id, formHabitName, hasDown)
	components.Habit(habit).Render(r.Context(), w)
}

func (h *HabitHandler) Modal(w http.ResponseWriter, r *http.Request) {
	components.HabitModal().Render(r.Context(), w)
}
