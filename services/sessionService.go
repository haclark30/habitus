package services

import (
	"fmt"
	"habitus/models"
	"log/slog"

	"github.com/google/uuid"
)

type SessionStore interface {
	GetSession(sessionToken string) (models.User, error)
	SaveSession(userId int, sessionKey string) error
}

type SessionService struct {
	Log          *slog.Logger
	SessionStore SessionStore
	UserStore    UserStore
}

func NewSessionService(log *slog.Logger, sessionStore SessionStore, userStore UserStore) SessionService {
	return SessionService{Log: log, SessionStore: sessionStore, UserStore: userStore}
}

func (s SessionService) GetSession(sessionToken string) (models.User, error) {
	return s.SessionStore.GetSession(sessionToken)
}

func (s SessionService) CreateSession(username string) (string, error) {
	user, err := s.UserStore.GetUser(username)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}
	sessionToken := uuid.NewString()
	err = s.SessionStore.SaveSession(user.Id, sessionToken)
	if err != nil {
		return "", fmt.Errorf("error saving session: %w", err)
	}
	return sessionToken, nil
}
