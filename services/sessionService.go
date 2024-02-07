package services

import (
	"context"
	"fmt"
	"habitus/db_sqlc"
	"log/slog"

	"github.com/google/uuid"
)

type SessionStore interface {
	GetSession(context.Context, string) (db_sqlc.User, error)
	AddSession(context.Context, db_sqlc.AddSessionParams) (db_sqlc.Session, error)
}

type SessionService struct {
	Log          *slog.Logger
	SessionStore SessionStore
	UserStore    UserStore
}

func NewSessionService(log *slog.Logger, sessionStore SessionStore, userStore UserStore) SessionService {
	return SessionService{Log: log, SessionStore: sessionStore, UserStore: userStore}
}

func (s SessionService) GetSession(sessionToken string) (db_sqlc.User, error) {
	return s.SessionStore.GetSession(context.TODO(), sessionToken)
}

func (s SessionService) CreateSession(username string) (string, error) {
	user, err := s.UserStore.GetUser(context.TODO(), username)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}
	sessionToken := uuid.NewString()
	_, err = s.SessionStore.AddSession(context.TODO(), db_sqlc.AddSessionParams{Userid: user.ID, Token: sessionToken})
	if err != nil {
		return "", fmt.Errorf("error saving session: %w", err)
	}
	return sessionToken, nil
}
