// middleware/auth.go
package middleware

import (
	"acore/database/redis"
	"log/slog"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session_token"
	sessionTTL        = 24 * time.Hour
	userIDHeader      = "X-User-ID"
)

func AuthRequired(next http.Handler) http.Handler {
	slog.Info("AuthRequired")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readSessionToken(r)
		if err != nil {
			redirectLogin(w, r)
			return
		}

		userID, err := checkRedis(token)
		if err != nil {
			redirectLogin(w, r)
			return
		}

		_ = refreshTTL(token, userID)
		r.Header.Set(userIDHeader, userID)
		slog.Info("Setting user id to header")

		next.ServeHTTP(w, r)
	})
}

func PublicOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := readSessionToken(r); err == nil {
			http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func readSessionToken(r *http.Request) (string, error) {
	slog.Info("Reading Cookie")
	c, err := r.Cookie(sessionCookieName)
	if err != nil || c.Value == "" {
		slog.Info("No Cookie Found")
		return "", err
	}
	return c.Value, nil
}

func checkRedis(token string) (string, error) {
	slog.Info("Checking redis")
	return redis.GetRedis(token)
}

func refreshTTL(token, userID string) error {
	return redis.SetRedis(token, userID, sessionTTL)
}

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("Redirect to login")
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
