package middleware

import (
	"context"
	"habitus/db_sqlc"
	"net/http"
)

type key int

const (
	userKey key = iota
)

type SessionService interface {
	GetSession(sessionToken string) (db_sqlc.User, error)
	CreateSession(username string) (string, error)
}

type SessionManager struct {
	sessionService SessionService
}

func NewSessionManager(sessionService SessionService) *SessionManager {
	return &SessionManager{
		sessionService: sessionService,
	}
}

func (s *SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("hx-request") == "true" {
				w.Header().Set("hx-redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionToken := c.Value
		user, err := s.sessionService.GetSession(sessionToken)
		if err != nil {
			if r.Header.Get("hx-request") == "true" {
				w.Header().Set("hx-redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) *db_sqlc.User {
	user := ctx.Value(userKey)
	if user == nil {
		return nil
	}
	userVal := user.(db_sqlc.User)
	return &userVal
}
