// middleware/auth.go
package middleware

import (
	"acore/database/redis"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session_token"
	sessionTTL        = 24 * time.Hour
	userIDHeader      = "X-User-ID"
)

func AuthRequired(next http.Handler) http.Handler {
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
	c, err := r.Cookie(sessionCookieName)
	if err != nil || c.Value == "" {
		return "", err
	}
	return c.Value, nil
}

func checkRedis(token string) (string, error) {
	return redis.GetRedis(token)
}

func refreshTTL(token, userID string) error {
	return redis.SetRedis(token, userID, sessionTTL)
}

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
