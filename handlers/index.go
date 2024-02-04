package handlers

import (
	"habitus/components"
	"habitus/middleware"
	"net/http"
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
	components.Page(i.habitService.GetHabits(user.Id), i.dailyService.GetDailys(user.Id)).Render(r.Context(), w)
}
