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

	habits := i.habitService.GetHabits(int(user.ID))
	components.Page(habits, i.dailyService.GetDailys(int(user.ID))).Render(r.Context(), w)
}
