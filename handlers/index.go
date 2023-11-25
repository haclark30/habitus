package handlers

import (
	"habitus/components"
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
	components.Page(i.habitService.GetHabits(), i.dailyService.GetDailys()).Render(r.Context(), w)
}
