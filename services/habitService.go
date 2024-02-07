package services

import (
	"context"
	"database/sql"
	"habitus/db_sqlc"
	"log/slog"
)

type Habit struct {
	Log    *slog.Logger
	Habits *db_sqlc.Queries
}

type HabitAndLog struct {
	Habit    db_sqlc.Habit
	HabitLog db_sqlc.HabitLog
}

func NewHabit(log *slog.Logger, queries *db_sqlc.Queries) Habit {
	return Habit{
		Log:    log,
		Habits: queries,
	}
}

func (hs Habit) CountUp(habitId int) HabitAndLog {
	habLog, err := hs.Habits.HabitLogUp(context.TODO(), db_sqlc.HabitLogUpParams{
		Habitid:  int64(habitId),
		Datetime: 0,
	})
	if err != nil {
		hs.Log.Error("error counting up", "err", err)
	}
	hs.Log.Info("got hablog", "hablog", habLog)

	habit, _ := hs.Habits.GetHabit(context.TODO(), int64(habitId))

	return HabitAndLog{Habit: habit, HabitLog: habLog}
}

func (hs Habit) CountDown(habitId int) HabitAndLog {
	habLog, _ := hs.Habits.HabitLogDown(context.TODO(), db_sqlc.HabitLogDownParams{
		Habitid:  int64(habitId),
		Datetime: 0,
	})

	habit, _ := hs.Habits.GetHabit(context.TODO(), int64(habitId))

	return HabitAndLog{Habit: habit, HabitLog: habLog}
}

func (hs Habit) AddHabit(userId int, habitName string, hasUp, hasDown bool) HabitAndLog {
	hs.Log.Info("adding habit", "habitName", habitName)
	habit, err := hs.Habits.AddHabit(context.TODO(), db_sqlc.AddHabitParams{
		Userid:  int64(userId),
		Name:    habitName,
		Hasup:   hasUp,
		Hasdown: hasDown,
	})
	if err != nil {
		hs.Log.Error("error adding habit", "err", err)
	}

	habLog, err := hs.Habits.AddHabitLog(context.TODO(), db_sqlc.AddHabitLogParams{
		Habitid:  habit.ID,
		Datetime: 0,
	})
	if err != nil {
		hs.Log.Error("error adding habit log", "err", err)
	}
	return HabitAndLog{Habit: habit, HabitLog: habLog}
}

func (hs Habit) GetHabits(userId int) []HabitAndLog {
	var habsAndLogs []HabitAndLog
	habit, err := hs.Habits.GetHabits(context.TODO(), int64(userId))
	if err != nil {
		hs.Log.Error("error getting habits", "err", err)
	}
	for _, hab := range habit {
		habLog, err := hs.Habits.GetHabitLog(context.TODO(), db_sqlc.GetHabitLogParams{
			Habitid:  hab.ID,
			Datetime: 0,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				habLog, _ = hs.Habits.AddHabitLog(context.TODO(), db_sqlc.AddHabitLogParams{
					Habitid:  hab.ID,
					Datetime: 0,
				})
			}
		}
		habsAndLogs = append(habsAndLogs, HabitAndLog{Habit: hab, HabitLog: habLog})

	}
	return habsAndLogs
}
