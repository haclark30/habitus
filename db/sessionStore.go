package db

import (
	"database/sql"
	"fmt"
	"habitus/models"
	"log/slog"
)

type SessionStore struct {
	Log *slog.Logger
	db  *sql.DB
}

func NewSessionStore(log *slog.Logger, dbase *sql.DB) SessionStore {
	return SessionStore{Log: log, db: dbase}
}

func (s SessionStore) GetSession(sessionToken string) (models.User, error) {
	user := models.User{}
	var userId int

	tx, err := s.db.Begin()
	if err != nil {
		return user, fmt.Errorf("error starting txn: %w", err)
	}
	sessionStmt, err := tx.Prepare(`SELECT userId FROM sessions WHERE token = ?`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer sessionStmt.Close()
	err = sessionStmt.QueryRow(sessionToken).Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		return user, fmt.Errorf("error getting session: %w", err)
	}

	userStmt, err := tx.Prepare(`SELECT * FROM users WHERE id = ?`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer userStmt.Close()
	err = userStmt.QueryRow(userId).Scan(&user.Id, &user.UserName, &user.PasswordHash)
	if err != nil {
		_ = tx.Rollback()
		return user, fmt.Errorf("error getting user: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return user, fmt.Errorf("error commiting txn: %w", err)
	}
	return user, nil
}

func (s SessionStore) SaveSession(userId int, sessionKey string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting txn: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO sessions (token, userId) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(sessionKey, userId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error saving session: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error commiting txn: %w", err)
	}

	return nil
}
