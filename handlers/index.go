package handlers

import (
	"database/sql"
	"habitus/components"
	"habitus/db_sqlc"
	"habitus/middleware"
	"log/slog"
	"net/http"
	"time"
)

type IndexHandler struct {
	Log     *slog.Logger
	queries *db_sqlc.Queries
}

func NewIndexHandler(log *slog.Logger, queries *db_sqlc.Queries) *IndexHandler {
	return &IndexHandler{Log: log, queries: queries}
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

	// todo: make this faster?
	habitsAndLogs, err := i.queries.GetHabitsAndLogs(r.Context(), db_sqlc.GetHabitsAndLogsParams{
		Userid:   user.ID,
		Datetime: t.Unix(),
	})

	if len(habitsAndLogs) == 0 {
		habits, _ := i.queries.GetHabits(r.Context(), user.ID)
		for _, h := range habits {
			_, err := i.queries.GetHabitLog(r.Context(), db_sqlc.GetHabitLogParams{
				Habitid:  h.ID,
				Datetime: t.Unix(),
			})
			if err == sql.ErrNoRows {
				_, err := i.queries.AddHabitLog(r.Context(), db_sqlc.AddHabitLogParams{
					Habitid:  h.ID,
					Datetime: t.Unix(),
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}

	habitsAndLogs, err = i.queries.GetHabitsAndLogs(r.Context(), db_sqlc.GetHabitsAndLogsParams{
		Userid:   user.ID,
		Datetime: t.Unix(),
	})

	if err != nil {
		i.Log.Error("error getting habits", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo: make this faster
	dailysAndLogs, err := i.queries.GetDailysAndLogs(r.Context(), db_sqlc.GetDailysAndLogsParams{
		Userid:   user.ID,
		Datetime: t.Unix(),
	})

	if len(dailysAndLogs) == 0 {
		dailys, _ := i.queries.GetDailys(r.Context(), user.ID)
		for _, d := range dailys {
			_, err := i.queries.GetDailyLog(r.Context(), db_sqlc.GetDailyLogParams{
				Dailyid:  d.ID,
				Datetime: t.Unix(),
			})
			if err == sql.ErrNoRows {
				_, err := i.queries.AddDailyLog(r.Context(), db_sqlc.AddDailyLogParams{
					Dailyid:  d.ID,
					Datetime: t.Unix(),
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}

	dailysAndLogs, err = i.queries.GetDailysAndLogs(r.Context(), db_sqlc.GetDailysAndLogsParams{
		Userid:   user.ID,
		Datetime: t.Unix(),
	})

	if err != nil {
		i.Log.Error("error getting dailys", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	components.Page(habitsAndLogs, dailysAndLogs).Render(r.Context(), w)
}
