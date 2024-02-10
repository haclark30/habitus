package handlers

import (
	"fmt"
	"habitus/components"
	"habitus/db_sqlc"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Log     *slog.Logger
	queries *db_sqlc.Queries
}

type UserService interface {
	CreateUser(username, password string) error
	AuthUser(username, password string) (bool, error)
}

type SessionService interface {
	GetSession(sessionToken string) (db_sqlc.User, error)
	CreateSession(username string) (string, error)
}

func NewUserHandler(log *slog.Logger, queries *db_sqlc.Queries) *UserHandler {
	return &UserHandler{Log: log, queries: queries}
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

	dbUser, err := u.queries.GetUser(r.Context(), user)
	if err != nil {
		u.Log.Error("error trying to auth", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Passwordhash), []byte(pass))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	_, err = u.queries.AddSession(r.Context(), db_sqlc.AddSessionParams{
		Userid: dbUser.ID,
		Token:  sessionToken,
	})
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Log.Error("error creating user", "err", err)
		return
	}

	if err = u.queries.AddUser(r.Context(), db_sqlc.AddUserParams{
		Username: user, Passwordhash: string(passwordHash),
	}); err != nil {
		u.Log.Error("error creating user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
