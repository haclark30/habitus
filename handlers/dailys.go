package handlers

import (
	"habitus/components"
	"habitus/db_sqlc"
	"habitus/middleware"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type DailyHandler struct {
	Log     *slog.Logger
	queries *db_sqlc.Queries
}

type DailyService interface {
	CompleteDaily(dailyLogId int) db_sqlc.GetDailysAndLogsRow
	AddDaily(userId int, dailyName string) (db_sqlc.GetDailysAndLogsRow, error)
	GetDailys(userId int, dateTime int64) []db_sqlc.GetDailysAndLogsRow
}

func NewDailyHandler(log *slog.Logger, queries *db_sqlc.Queries) *DailyHandler {
	return &DailyHandler{Log: log, queries: queries}
}

func (d *DailyHandler) CompleteDaily(w http.ResponseWriter, r *http.Request) {
	dailyId, _ := strconv.Atoi(chi.URLParam(r, "dailyId"))
	dailyLogId, _ := strconv.Atoi(chi.URLParam(r, "dailyLogId"))
	dailyLog, err := d.queries.DailyLogDone(r.Context(), int64(dailyLogId))
	if err != nil {
		d.Log.Error("error marking daily done", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	daily, err := d.queries.GetDaily(r.Context(), int64(dailyId))
	if err != nil {
		d.Log.Error("error getting daily", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	components.Daily(daily, dailyLog).Render(r.Context(), w)
}

func (d *DailyHandler) Put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formDailyName := r.Form.Get("dailyName")
	user := middleware.UserFromContext(r.Context())
	daily, _ := d.queries.AddDaily(r.Context(), db_sqlc.AddDailyParams{
		Userid: user.ID,
		Name:   formDailyName,
	})

	loc, _ := time.LoadLocation("America/Chicago")
	year, month, day := time.Now().In(loc).Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, loc)

	dailyLog, _ := d.queries.AddDailyLog(r.Context(), db_sqlc.AddDailyLogParams{
		Dailyid:  daily.ID,
		Datetime: t.Unix(),
	})
	components.Daily(daily, dailyLog).Render(r.Context(), w)
}

func (d *DailyHandler) Modal(w http.ResponseWriter, r *http.Request) {
	components.DailyModal().Render(r.Context(), w)
}
