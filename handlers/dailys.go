package handlers

import (
	"habitus/components"
	"habitus/middleware"
	"habitus/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DailyHandler struct {
	Log          *slog.Logger
	DailyService DailyService
}

type DailyService interface {
	CompleteDaily(dailyId int) models.Daily
	AddDaily(userId int, dailyName string) models.Daily
	GetDailys(userId int) []models.Daily
}

func NewDailyHandler(log *slog.Logger, dailyService DailyService) *DailyHandler {
	return &DailyHandler{Log: log, DailyService: dailyService}
}

func (d *DailyHandler) CompleteDaily(w http.ResponseWriter, r *http.Request) {
	dailyIndex, _ := strconv.Atoi(chi.URLParam(r, "dailyId"))
	daily := d.DailyService.CompleteDaily(dailyIndex)
	components.Daily(daily).Render(r.Context(), w)
}

func (d *DailyHandler) Put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formDailyName := r.Form.Get("dailyName")
	user := middleware.UserFromContext(r.Context())
	daily := d.DailyService.AddDaily(int(user.ID), formDailyName)
	components.Daily(daily).Render(r.Context(), w)
}

func (d *DailyHandler) Modal(w http.ResponseWriter, r *http.Request) {
	components.DailyModal().Render(r.Context(), w)
}
