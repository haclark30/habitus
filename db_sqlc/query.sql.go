// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db_sqlc

import (
	"context"
)

const addDaily = `-- name: AddDaily :one
INSERT INTO dailys (
  userId, name 
) VALUES (
  ?, ?
) RETURNING id, userid, name
`

type AddDailyParams struct {
	Userid int64
	Name   string
}

func (q *Queries) AddDaily(ctx context.Context, arg AddDailyParams) (Daily, error) {
	row := q.db.QueryRowContext(ctx, addDaily, arg.Userid, arg.Name)
	var i Daily
	err := row.Scan(&i.ID, &i.Userid, &i.Name)
	return i, err
}

const addDailyLog = `-- name: AddDailyLog :one
INSERT INTO dailyLog (
  dailyId, done, dateTime
) VALUES (
  ?, false, ?
) RETURNING id, dailyid, done, datetime
`

type AddDailyLogParams struct {
	Dailyid  int64
	Datetime int64
}

func (q *Queries) AddDailyLog(ctx context.Context, arg AddDailyLogParams) (DailyLog, error) {
	row := q.db.QueryRowContext(ctx, addDailyLog, arg.Dailyid, arg.Datetime)
	var i DailyLog
	err := row.Scan(
		&i.ID,
		&i.Dailyid,
		&i.Done,
		&i.Datetime,
	)
	return i, err
}

const addHabit = `-- name: AddHabit :one
INSERT INTO habits (
  userId, name, hasUp, hasDown
) VALUES (
  ?, ?, ?, ?
)
RETURNING id, userid, name, hasup, hasdown
`

type AddHabitParams struct {
	Userid  int64
	Name    string
	Hasup   bool
	Hasdown bool
}

func (q *Queries) AddHabit(ctx context.Context, arg AddHabitParams) (Habit, error) {
	row := q.db.QueryRowContext(ctx, addHabit,
		arg.Userid,
		arg.Name,
		arg.Hasup,
		arg.Hasdown,
	)
	var i Habit
	err := row.Scan(
		&i.ID,
		&i.Userid,
		&i.Name,
		&i.Hasup,
		&i.Hasdown,
	)
	return i, err
}

const addHabitLog = `-- name: AddHabitLog :one
INSERT INTO habitLog (
  habitId, upCount, downCount, dateTime
) VALUES (
  ?, 0, 0, ?
)
RETURNING id, habitid, upcount, downcount, datetime
`

type AddHabitLogParams struct {
	Habitid  int64
	Datetime int64
}

func (q *Queries) AddHabitLog(ctx context.Context, arg AddHabitLogParams) (HabitLog, error) {
	row := q.db.QueryRowContext(ctx, addHabitLog, arg.Habitid, arg.Datetime)
	var i HabitLog
	err := row.Scan(
		&i.ID,
		&i.Habitid,
		&i.Upcount,
		&i.Downcount,
		&i.Datetime,
	)
	return i, err
}

const addSession = `-- name: AddSession :one
INSERT INTO sessions (
  token, userId
) VALUES (
  ?, ?
)
RETURNING token, userid
`

type AddSessionParams struct {
	Token  string
	Userid int64
}

func (q *Queries) AddSession(ctx context.Context, arg AddSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, addSession, arg.Token, arg.Userid)
	var i Session
	err := row.Scan(&i.Token, &i.Userid)
	return i, err
}

const addUser = `-- name: AddUser :exec
INSERT INTO users (
  username, passwordHash
) VALUES (
  ?, ?
)
`

type AddUserParams struct {
	Username     string
	Passwordhash string
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) error {
	_, err := q.db.ExecContext(ctx, addUser, arg.Username, arg.Passwordhash)
	return err
}

const dailyLogDone = `-- name: DailyLogDone :one
UPDATE dailyLog
SET done = NOT done
WHERE ID = ?
RETURNING id, dailyid, done, datetime
`

func (q *Queries) DailyLogDone(ctx context.Context, id int64) (DailyLog, error) {
	row := q.db.QueryRowContext(ctx, dailyLogDone, id)
	var i DailyLog
	err := row.Scan(
		&i.ID,
		&i.Dailyid,
		&i.Done,
		&i.Datetime,
	)
	return i, err
}

const deleteHabit = `-- name: DeleteHabit :exec
DELETE FROM habits
WHERE ID = ?
`

func (q *Queries) DeleteHabit(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteHabit, id)
	return err
}

const deleteHabitLogs = `-- name: DeleteHabitLogs :exec
DELETE FROM habitLog
WHERE habitId = ?
`

func (q *Queries) DeleteHabitLogs(ctx context.Context, habitid int64) error {
	_, err := q.db.ExecContext(ctx, deleteHabitLogs, habitid)
	return err
}

const getDaily = `-- name: GetDaily :one
SELECT id, userid, name FROM dailys
WHERE ID = ?
`

func (q *Queries) GetDaily(ctx context.Context, id int64) (Daily, error) {
	row := q.db.QueryRowContext(ctx, getDaily, id)
	var i Daily
	err := row.Scan(&i.ID, &i.Userid, &i.Name)
	return i, err
}

const getDailyLog = `-- name: GetDailyLog :one
SELECT id, dailyid, done, datetime from dailyLog
WHERE dailyId = ? and dateTime = ?
LIMIT 1
`

type GetDailyLogParams struct {
	Dailyid  int64
	Datetime int64
}

func (q *Queries) GetDailyLog(ctx context.Context, arg GetDailyLogParams) (DailyLog, error) {
	row := q.db.QueryRowContext(ctx, getDailyLog, arg.Dailyid, arg.Datetime)
	var i DailyLog
	err := row.Scan(
		&i.ID,
		&i.Dailyid,
		&i.Done,
		&i.Datetime,
	)
	return i, err
}

const getDailys = `-- name: GetDailys :many
SELECT id, userid, name FROM dailys
WHERE userId = ?
`

func (q *Queries) GetDailys(ctx context.Context, userid int64) ([]Daily, error) {
	rows, err := q.db.QueryContext(ctx, getDailys, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Daily
	for rows.Next() {
		var i Daily
		if err := rows.Scan(&i.ID, &i.Userid, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDailysAndLogs = `-- name: GetDailysAndLogs :many
SELECT d.id, d.userid, d.name, dl.id, dl.dailyid, dl.done, dl.datetime FROM dailys d
JOIN dailyLog dl ON dl.dailyId = d.ID
WHERE d.userId = ? and dl.dateTime = ?
`

type GetDailysAndLogsParams struct {
	Userid   int64
	Datetime int64
}

type GetDailysAndLogsRow struct {
	Daily    Daily
	DailyLog DailyLog
}

func (q *Queries) GetDailysAndLogs(ctx context.Context, arg GetDailysAndLogsParams) ([]GetDailysAndLogsRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailysAndLogs, arg.Userid, arg.Datetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDailysAndLogsRow
	for rows.Next() {
		var i GetDailysAndLogsRow
		if err := rows.Scan(
			&i.Daily.ID,
			&i.Daily.Userid,
			&i.Daily.Name,
			&i.DailyLog.ID,
			&i.DailyLog.Dailyid,
			&i.DailyLog.Done,
			&i.DailyLog.Datetime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHabit = `-- name: GetHabit :one
SELECT id, userid, name, hasup, hasdown FROM habits
WHERE id = ? LIMIT 1
`

func (q *Queries) GetHabit(ctx context.Context, id int64) (Habit, error) {
	row := q.db.QueryRowContext(ctx, getHabit, id)
	var i Habit
	err := row.Scan(
		&i.ID,
		&i.Userid,
		&i.Name,
		&i.Hasup,
		&i.Hasdown,
	)
	return i, err
}

const getHabitLog = `-- name: GetHabitLog :one
SELECT id, habitid, upcount, downcount, datetime from habitLog
WHERE habitId = ? and dateTime = ?
LIMIT 1
`

type GetHabitLogParams struct {
	Habitid  int64
	Datetime int64
}

func (q *Queries) GetHabitLog(ctx context.Context, arg GetHabitLogParams) (HabitLog, error) {
	row := q.db.QueryRowContext(ctx, getHabitLog, arg.Habitid, arg.Datetime)
	var i HabitLog
	err := row.Scan(
		&i.ID,
		&i.Habitid,
		&i.Upcount,
		&i.Downcount,
		&i.Datetime,
	)
	return i, err
}

const getHabits = `-- name: GetHabits :many
SELECT id, userid, name, hasup, hasdown FROM habits 
WHERE userId = ?
`

func (q *Queries) GetHabits(ctx context.Context, userid int64) ([]Habit, error) {
	rows, err := q.db.QueryContext(ctx, getHabits, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Habit
	for rows.Next() {
		var i Habit
		if err := rows.Scan(
			&i.ID,
			&i.Userid,
			&i.Name,
			&i.Hasup,
			&i.Hasdown,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHabitsAndLogs = `-- name: GetHabitsAndLogs :many
SELECT habits.id, habits.userid, habits.name, habits.hasup, habits.hasdown, hl.id, hl.habitid, hl.upcount, hl.downcount, hl.datetime FROM habits 
JOIN habitLog hl ON hl.habitId = habits.ID
WHERE habits.userId = ? and hl.dateTime = ?
`

type GetHabitsAndLogsParams struct {
	Userid   int64
	Datetime int64
}

type GetHabitsAndLogsRow struct {
	Habit    Habit
	HabitLog HabitLog
}

func (q *Queries) GetHabitsAndLogs(ctx context.Context, arg GetHabitsAndLogsParams) ([]GetHabitsAndLogsRow, error) {
	rows, err := q.db.QueryContext(ctx, getHabitsAndLogs, arg.Userid, arg.Datetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetHabitsAndLogsRow
	for rows.Next() {
		var i GetHabitsAndLogsRow
		if err := rows.Scan(
			&i.Habit.ID,
			&i.Habit.Userid,
			&i.Habit.Name,
			&i.Habit.Hasup,
			&i.Habit.Hasdown,
			&i.HabitLog.ID,
			&i.HabitLog.Habitid,
			&i.HabitLog.Upcount,
			&i.HabitLog.Downcount,
			&i.HabitLog.Datetime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSession = `-- name: GetSession :one
SELECT users.id, users.username, users.passwordhash FROM sessions
JOIN users ON users.ID = sessions.userId 
WHERE token = ?
`

func (q *Queries) GetSession(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getSession, token)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Passwordhash)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, passwordhash FROM users
WHERE username = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Passwordhash)
	return i, err
}

const habitLogDown = `-- name: HabitLogDown :one
UPDATE habitLog
SET downCount = downCount + 1
WHERE ID = ? 
RETURNING id, habitid, upcount, downcount, datetime
`

func (q *Queries) HabitLogDown(ctx context.Context, id int64) (HabitLog, error) {
	row := q.db.QueryRowContext(ctx, habitLogDown, id)
	var i HabitLog
	err := row.Scan(
		&i.ID,
		&i.Habitid,
		&i.Upcount,
		&i.Downcount,
		&i.Datetime,
	)
	return i, err
}

const habitLogUp = `-- name: HabitLogUp :one
UPDATE habitLog
SET upCount = upCount + 1
WHERE ID = ?
RETURNING id, habitid, upcount, downcount, datetime
`

func (q *Queries) HabitLogUp(ctx context.Context, id int64) (HabitLog, error) {
	row := q.db.QueryRowContext(ctx, habitLogUp, id)
	var i HabitLog
	err := row.Scan(
		&i.ID,
		&i.Habitid,
		&i.Upcount,
		&i.Downcount,
		&i.Datetime,
	)
	return i, err
}

const updateHabit = `-- name: UpdateHabit :one
UPDATE habits 
SET 
  name = ?,
  hasUp = ?,
  hasDown = ?
WHERE
  id = ?
RETURNING id, userid, name, hasup, hasdown
`

type UpdateHabitParams struct {
	Name    string
	Hasup   bool
	Hasdown bool
	ID      int64
}

func (q *Queries) UpdateHabit(ctx context.Context, arg UpdateHabitParams) (Habit, error) {
	row := q.db.QueryRowContext(ctx, updateHabit,
		arg.Name,
		arg.Hasup,
		arg.Hasdown,
		arg.ID,
	)
	var i Habit
	err := row.Scan(
		&i.ID,
		&i.Userid,
		&i.Name,
		&i.Hasup,
		&i.Hasdown,
	)
	return i, err
}

const updateHabitLog = `-- name: UpdateHabitLog :one
UPDATE habitLog
SET 
  upCount = ?,
  downCount = ?
WHERE
  id = ?
RETURNING id, habitid, upcount, downcount, datetime
`

type UpdateHabitLogParams struct {
	Upcount   int64
	Downcount int64
	ID        int64
}

func (q *Queries) UpdateHabitLog(ctx context.Context, arg UpdateHabitLogParams) (HabitLog, error) {
	row := q.db.QueryRowContext(ctx, updateHabitLog, arg.Upcount, arg.Downcount, arg.ID)
	var i HabitLog
	err := row.Scan(
		&i.ID,
		&i.Habitid,
		&i.Upcount,
		&i.Downcount,
		&i.Datetime,
	)
	return i, err
}
