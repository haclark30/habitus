-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: AddUser :exec
INSERT INTO users (
  username, passwordHash
) VALUES (
  ?, ?
);

-- name: AddSession :one
INSERT INTO sessions (
  token, userId
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetSession :one
SELECT users.* FROM sessions
JOIN users ON users.ID = sessions.userId 
WHERE token = ?;

-- name: GetHabits :many
SELECT * FROM habits 
WHERE userId = ?;

-- name: GetHabit :one
SELECT * FROM habits
WHERE id = ? LIMIT 1;

-- name: AddHabit :one
INSERT INTO habits (
  userId, name, hasUp, hasDown
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: AddHabitLog :one
INSERT INTO habitLog (
  habitId, upCount, downCount, dateTime
) VALUES (
  ?, 0, 0, ?
)
RETURNING *;

-- name: HabitLogUp :one
UPDATE habitLog
SET upCount = upCount + 1
WHERE ID = ?
RETURNING *;

-- name: HabitLogDown :one
UPDATE habitLog
SET downCount = downCount + 1
WHERE ID = ? 
RETURNING *;

-- name: GetHabitLog :one
SELECT * from habitLog
WHERE habitId = ? and dateTime = ?
LIMIT 1;

-- name: GetHabitsAndLogs :many
SELECT sqlc.embed(habits), sqlc.embed(hl) FROM habits 
JOIN habitLog hl ON hl.habitId = habits.ID
WHERE habits.userId = ?;
