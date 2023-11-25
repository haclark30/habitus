package db

import "habitus/models"

type HabitStore struct {
	habits *[]models.Habit
}

type DailyStore struct {
	dailys *[]models.Daily
}

func NewHabitStore(habits *[]models.Habit) HabitStore {
	return HabitStore{habits: habits}
}

func NewDailyStore(dailys *[]models.Daily) DailyStore {
	return DailyStore{dailys: dailys}
}

func (hs HabitStore) AddHabit(habitName string, hasDown bool) models.Habit {
	habit := models.Habit{
		Id:      len(*(hs.habits)),
		Name:    habitName,
		HasDown: hasDown,
	}
	*(hs.habits) = append(*(hs.habits), habit)
	return habit
}

func (hs HabitStore) GetHabits() []models.Habit {
	return *hs.habits
}

func (hs HabitStore) GetHabit(habitId int) models.Habit {
	return (*hs.habits)[habitId]
}

func (hs HabitStore) UpdateHabit(habit models.Habit) models.Habit {
	(*hs.habits)[habit.Id] = habit
	return habit
}

func (ds DailyStore) AddDaily(dailyName string) models.Daily {
	daily := models.Daily{
		Id:   len(*(ds.dailys)),
		Name: dailyName,
	}
	*(ds.dailys) = append(*(ds.dailys), daily)
	return daily
}

func (ds DailyStore) UpdateDaily(daily models.Daily) models.Daily {
	(*ds.dailys)[daily.Id] = daily
	return daily
}

func (ds DailyStore) GetDailys() []models.Daily {
	return *ds.dailys
}

func (ds DailyStore) GetDaily(dailyId int) models.Daily {
	return (*ds.dailys)[dailyId]
}
