package handlers

import (
	"fmt"
	"habitus/components"
	"habitus/db_sqlc"
	"habitus/middleware"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type HabitHandler struct {
	Log     *slog.Logger
	queries *db_sqlc.Queries
}

func NewHabitHandler(log *slog.Logger, queries *db_sqlc.Queries) *HabitHandler {
	return &HabitHandler{Log: log, queries: queries}
}

func (h *HabitHandler) CountUp(w http.ResponseWriter, r *http.Request) {
	habitIndex, _ := strconv.Atoi(chi.URLParam(r, "habitId"))
	habitLogId, _ := strconv.Atoi(chi.URLParam(r, "habitLogId"))
	habitLog, err := h.queries.HabitLogUp(r.Context(), int64(habitLogId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	habit, err := h.queries.GetHabit(r.Context(), int64(habitIndex))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	components.Habit(habit, habitLog).Render(r.Context(), w)
}

func (h *HabitHandler) CountDown(w http.ResponseWriter, r *http.Request) {
	habitIndex, _ := strconv.Atoi(chi.URLParam(r, "habitId"))
	habitLogId, _ := strconv.Atoi(chi.URLParam(r, "habitLogId"))
	habitLog, err := h.queries.HabitLogDown(r.Context(), int64(habitLogId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	habit, err := h.queries.GetHabit(r.Context(), int64(habitIndex))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	components.Habit(habit, habitLog).Render(r.Context(), w)
}

func (h *HabitHandler) Put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formHabitName := r.Form.Get("habitName")
	formHasUp := r.Form.Get("hasUp")
	formHasDown := r.Form.Get("hasDown")

	hasUp := formHasUp == "on"
	hasDown := formHasDown == "on"
	user := middleware.UserFromContext(r.Context())
	habit, err := h.queries.AddHabit(r.Context(), db_sqlc.AddHabitParams{
		Userid:  user.ID,
		Name:    formHabitName,
		Hasup:   hasUp,
		Hasdown: hasDown,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	loc, _ := time.LoadLocation("America/Chicago")
	year, month, day := time.Now().In(loc).Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, loc)

	habitLog, err := h.queries.AddHabitLog(r.Context(), db_sqlc.AddHabitLogParams{
		Habitid:  habit.ID,
		Datetime: t.Unix(),
	})
	components.Habit(habit, habitLog).Render(r.Context(), w)
}

func (h *HabitHandler) Modal(w http.ResponseWriter, r *http.Request) {
	components.HabitModal().Render(r.Context(), w)
}

func (h *HabitHandler) Edit(w http.ResponseWriter, r *http.Request) {
	habitId, _ := strconv.Atoi(chi.URLParam(r, "habitId"))
	habit, _ := h.queries.GetHabit(r.Context(), int64(habitId))

	fmt.Fprintf(w, "got %v", habit.Name)
}
