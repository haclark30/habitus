package services

import (
	"habitus/db"
	"habitus/models"
	"log/slog"
)

type Habit struct {
	Log    *slog.Logger
	Habits db.HabitStore
}

func NewHabit(log *slog.Logger, habits db.HabitStore) Habit {
	return Habit{
		Log:    log,
		Habits: habits,
	}
}

func (hs Habit) CountUp(habitId int) models.Habit {
	habit := hs.Habits.GetHabit(habitId)
	habit.UpCount++
	return hs.Habits.UpdateHabit(habit)
}

func (hs Habit) CountDown(habitId int) models.Habit {
	habit := hs.Habits.GetHabit(habitId)
	habit.DownCount++
	return hs.Habits.UpdateHabit(habit)
}

func (hs Habit) AddHabit(userId int, habitName string, hasDown bool) models.Habit {
	return hs.Habits.AddHabit(userId, habitName, hasDown)
}

func (hs Habit) GetHabits(userId int) []models.Habit {
	return hs.Habits.GetHabits(userId)
}
