package db

import (
	"database/sql"
	"fmt"
	"habitus/models"
	"log"
	"log/slog"
)

type HabitStore struct {
	dbase *sql.DB
}

type DailyStore struct {
	dbase *sql.DB
}

type UserStore struct {
	log   *slog.Logger
	dbase *sql.DB
}

func NewHabitStore(dbase *sql.DB) HabitStore {
	return HabitStore{dbase: dbase}
}

func NewDailyStore(log *slog.Logger, dbase *sql.DB) DailyStore {
	return DailyStore{dbase: dbase}
}

func NewUserStore(log *slog.Logger, dbase *sql.DB) UserStore {
	return UserStore{log: log, dbase: dbase}
}

func (us UserStore) GetUser(username string) (models.User, error) {
	user := models.User{}
	// tx, err := us.dbase.Begin()
	// if err != nil {
	// 	return user, fmt.Errorf("error starting txn: %w", err)
	// }
	// stmt, err := tx.Prepare(`SELECT * FROM users WHERE username = ?`)
	// if err != nil {
	// 	return user, fmt.Errorf("error preparing statement: %w", err)
	// }
	// defer stmt.Close()

	err := us.dbase.QueryRow(`SELECT * FROM users WHERE username = ?`, username).Scan(&user.Id, &user.UserName, &user.PasswordHash)
	if err != nil {
		// _ = tx.Rollback()
		return user, fmt.Errorf("error getting user: %w", err)
	}

	// err = tx.Commit()
	if err != nil {
		return user, fmt.Errorf("error commiting txn: %w", err)
	}
	return user, nil
}

func (us UserStore) AddUser(username, passwordHash string) error {
	tx, err := us.dbase.Begin()
	if err != nil {
		us.log.Error("error starting txn", "err", err)
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO users
	(username, passwordHash) VALUES(?, ?)`)
	if err != nil {
		us.log.Error("error preparing statement", "err", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, passwordHash)
	if err != nil {
		_ = tx.Rollback()
		us.log.Error("error inserting user", "err", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		us.log.Error("error commiting txn", "err", err)
		return err
	}
	return nil
}

func (hs HabitStore) AddHabit(userId int, habitName string, hasDown bool) models.Habit {
	habit := models.Habit{
		Name:    habitName,
		HasDown: hasDown,
	}
	tx, err := hs.dbase.Begin()
	if err != nil {
		log.Fatal("error starting txn", "err", err)
	}
	stmt, err := tx.Prepare(`INSERT INTO habits
		(userId, name, up, down, hasDown) VALUES(?, ?, ?, ?, ?)
		RETURNING *`)
	if err != nil {
		log.Fatal("error preparing statement", "err", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(userId, habitName, 0, 0, hasDown).Scan(&habit.Id, &habit.UserId, &habit.Name, &habit.UpCount, &habit.DownCount, &habit.HasDown)

	if err != nil {
		log.Fatal("error querying row", "err", err)
	}
	err = tx.Commit()

	if err != nil {
		log.Fatal("error commiting txn", "err", err)
	}

	return habit
}

func (hs HabitStore) GetHabits(userId int) []models.Habit {
	// todo: handle error
	rows, err := hs.dbase.Query("SELECT * FROM HABITS WHERE userId = ?", userId)
	if err != nil {
		log.Fatal("error querying db", "err", err)
	}
	var habits []models.Habit
	for rows.Next() {
		var habit models.Habit
		// todo: handle error
		rows.Scan(&habit.Id, &habit.UserId, &habit.Name, &habit.UpCount, &habit.DownCount, &habit.HasDown)
		habits = append(habits, habit)
	}
	return habits
}

func (hs HabitStore) GetHabit(habitId int) models.Habit {
	rows, err := hs.dbase.Query("SELECT * FROM HABITS WHERE id = ?", habitId)
	if err != nil {
		log.Fatal("error querying db", "err", err)
	}
	var habits []models.Habit
	for rows.Next() {
		var habit models.Habit
		// todo: handle error
		rows.Scan(&habit.Id, &habit.UserId, &habit.Name, &habit.UpCount, &habit.DownCount, &habit.HasDown)
		habits = append(habits, habit)
	}
	return habits[0]
}

func (hs HabitStore) UpdateHabit(habit models.Habit) models.Habit {
	_, err := hs.dbase.Exec("UPDATE habits SET name=?, up=?, down=?, hasDown=? WHERE id = ?", habit.Name, habit.UpCount, habit.DownCount, habit.HasDown, habit.Id)
	if err != nil {
		log.Fatal(err)
	}
	return habit
}

func (ds DailyStore) AddDaily(userId int, dailyName string) models.Daily {
	daily := models.Daily{
		Name: dailyName,
	}
	tx, err := ds.dbase.Begin()
	if err != nil {
		log.Fatal("error starting txn", "err", err)
	}
	stmt, err := tx.Prepare(`INSERT INTO dailys(userId, name, due, done)
		VALUES (?, ?, ?, ?) RETURNING *`)
	if err != nil {
		log.Fatal("error preparing statement", "err", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(userId, dailyName, false, false).Scan(&daily.Id, &daily.UserId, &daily.Name, &daily.Due, &daily.Done)
	if err != nil {
		log.Fatal("error querying row", "err", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("error commiting txn", "err", err)
	}
	return daily
}

func (ds DailyStore) UpdateDaily(daily models.Daily) models.Daily {
	tx, err := ds.dbase.Begin()
	if err != nil {
		log.Fatal("error starting txn", "err", err)
	}
	stmt, err := tx.Prepare(`UPDATE dailys SET name=?, due=?, done=? WHERE id = ? RETURNING *`)
	if err != nil {
		log.Fatal("error preparing statement", "err", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(daily.Name, daily.Due, daily.Done, daily.Id).Scan(&daily.Id, &daily.UserId, &daily.Name, &daily.Due, &daily.Done)
	if err != nil {
		log.Fatal("error querying row", "err", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("error commiting txn", "err", err)
	}
	return daily
}

func (ds DailyStore) GetDailys(userId int) []models.Daily {
	rows, err := ds.dbase.Query("SELECT * FROM dailys WHERE userId = ?", userId)
	if err != nil {
		log.Fatal("error querying db", "err", err)
	}
	var dailys []models.Daily
	for rows.Next() {
		var daily models.Daily
		err := rows.Scan(&daily.Id, &daily.UserId, &daily.Name, &daily.Due, &daily.Done)
		if err != nil {
			log.Fatal("error scanning row", "err", err)
		}
		dailys = append(dailys, daily)
	}
	return dailys
}

func (ds DailyStore) GetDaily(dailyId int) models.Daily {
	var daily models.Daily
	err := ds.dbase.QueryRow(`SELECT * FROM dailys WHERE ID = ?`, dailyId).
		Scan(&daily.Id, &daily.UserId, &daily.Name, &daily.Due, &daily.Done)
	if err != nil {
		log.Fatal("error querying db", "err", err)
	}
	return daily
}
