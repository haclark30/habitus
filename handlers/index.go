package handlers

import (
	"habitus/components"
	"habitus/middleware"
	"net/http"
	"time"
)

type IndexHandler struct {
	habitService HabitService
	dailyService DailyService
}

func NewIndexHandler(habitService HabitService, dailyService DailyService) *IndexHandler {
	return &IndexHandler{habitService: habitService, dailyService: dailyService}
}

func (i *IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	loc, _ := time.LoadLocation("America/Chicago")
	year, month, day := time.Now().In(loc).Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, loc)

	habits := i.habitService.GetHabits(int(user.ID), t.Unix())
	components.Page(habits, i.dailyService.GetDailys(int(user.ID))).Render(r.Context(), w)
}
