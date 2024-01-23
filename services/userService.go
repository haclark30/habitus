package services

import (
	"fmt"
	"habitus/models"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Log       *slog.Logger
	UserStore UserStore
}

type UserStore interface {
	GetUser(userName string) (models.User, error)
	AddUser(username, passwordHash string) error
}

func NewUserService(log *slog.Logger, userStore UserStore) *User {
	return &User{Log: log, UserStore: userStore}
}

func (u User) CreateUser(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		u.Log.Error("error hashing password", "err", err)
		return fmt.Errorf("error hashing password: %w", err)
	}
	if err = u.UserStore.AddUser(username, string(passwordHash)); err != nil {
		slog.Error("error adding user", "err", err)
		return fmt.Errorf("error adding user: %w", err)
	}
	return nil
}

func (u User) AuthUser(username, password string) (bool, error) {
	user, err := u.UserStore.GetUser(username)
	if err != nil {
		return false, fmt.Errorf("error getting user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
