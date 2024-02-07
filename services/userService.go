package services

import (
	"context"
	"fmt"
	"habitus/db_sqlc"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Log       *slog.Logger
	UserStore UserStore
}

type UserStore interface {
	GetUser(context.Context, string) (db_sqlc.User, error)
	AddUser(context.Context, db_sqlc.AddUserParams) error
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
	if err = u.UserStore.AddUser(context.TODO(), db_sqlc.AddUserParams{Username: username, Passwordhash: string(passwordHash)}); err != nil {
		slog.Error("error adding user", "err", err)
		return fmt.Errorf("error adding user: %w", err)
	}
	return nil
}

func (u User) AuthUser(username, password string) (bool, error) {
	user, err := u.UserStore.GetUser(context.TODO(), username)
	if err != nil {
		return false, fmt.Errorf("error getting user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Passwordhash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
