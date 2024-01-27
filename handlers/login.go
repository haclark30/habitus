package handlers

import (
	"fmt"
	"habitus/components"
	"habitus/models"
	"log/slog"
	"net/http"
	"time"
)

type UserHandler struct {
	Log            *slog.Logger
	UserService    UserService
	SessionService SessionService
}

type UserService interface {
	CreateUser(username, password string) error
	AuthUser(username, password string) (bool, error)
}

type SessionService interface {
	GetSession(sessionToken string) (models.User, error)
	CreateSession(username string) (string, error)
}

func NewUserHandler(log *slog.Logger, userService UserService, sessionService SessionService) *UserHandler {
	return &UserHandler{Log: log, UserService: userService, SessionService: sessionService}
}

func (u *UserHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	components.Login().Render(r.Context(), w)
}

func (u *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error parsing form", err)
	}
	user, pass := r.Form.Get("userName"), r.Form.Get("pass")
	authOk, err := u.UserService.AuthUser(user, pass)
	if err != nil {
		u.Log.Error("error trying to auth", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken, err := u.SessionService.CreateSession(user)
	if err != nil {
		u.Log.Error("error writing session", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expiresAt := time.Now().Add(24 * 365 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
	})

	w.Header().Set("hx-redirect", "/")
}

func (u *UserHandler) GetSignup(w http.ResponseWriter, r *http.Request) {
	components.Signup().Render(r.Context(), w)
}

func (u *UserHandler) PostSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Log.Error("error parsing form", "err", err)
		return
	}
	user, pass := r.Form.Get("userName"), r.Form.Get("pass")
	if user == "" || pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		u.Log.Error("user or password are blank")
	}
	err = u.UserService.CreateUser(user, pass)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Log.Error("error creating user", "err", err)
	}
	w.Header().Set("hx-redirect", "/")
}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})

	if w.Header().Get("hx-request") == "true" {
		w.Header().Add("hx-redirect", "/login")
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
