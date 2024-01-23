package handlers

import (
	"habitus/components"
	"habitus/models"
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
	userCtx := r.Context().Value("user")
	if userCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := userCtx.(models.User)
	components.Page(i.habitService.GetHabits(user.Id), i.dailyService.GetDailys(user.Id)).Render(r.Context(), w)
}
