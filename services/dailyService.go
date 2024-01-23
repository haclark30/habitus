package services

import (
	"habitus/db"
	"habitus/models"
	"log/slog"
)

type Daily struct {
	Log    *slog.Logger
	Dailys db.DailyStore
}

func NewDaily(log *slog.Logger, dailys db.DailyStore) Daily {
	return Daily{
		Log:    log,
		Dailys: dailys,
	}
}

func (ds Daily) AddDaily(userId int, dailyName string) models.Daily {
	return ds.Dailys.AddDaily(userId, dailyName)
}

func (ds Daily) CompleteDaily(dailyId int) models.Daily {
	daily := ds.Dailys.GetDaily(dailyId)
	daily.Done = !daily.Done
	ds.Dailys.UpdateDaily(daily)
	return daily
}

func (ds Daily) GetDailys(userId int) []models.Daily {
	return ds.Dailys.GetDailys(userId)
}
