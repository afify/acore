package middleware

import (
	"log/slog"
	"net/http"

	"acore/models/session"
)

const (
	sessionCookieName = "session_token"
	userIDHeader      = "X-User-ID"
)

func AuthRequired(next http.Handler) http.Handler {
	slog.Info("AuthRequired")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("AuthRequired: incoming request",
			slog.String("path", r.URL.Path),
		)

		// 1) read cookie
		c, err := r.Cookie(sessionCookieName)
		if err != nil || c.Value == "" {
			session.RedirectLogin(w, r)
			return
		}

		// 2) verify token (decrypt & check expiry)
		userID, err := session.VerifySessionToken(c.Value)
		if err != nil {
			session.RedirectLogin(w, r)
			return
		}

		// 3) inject userID and proceed
		r.Header.Set(userIDHeader, userID)
		slog.Info("Session valid, setting user ID header", slog.String("userID", userID))
		next.ServeHTTP(w, r)
	})
}

func PublicOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("PublicOnly: incoming request",
			slog.String("path", r.URL.Path),
		)
		if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
			if userID, err := session.VerifySessionToken(c.Value); err == nil {
				slog.Info("PublicOnly: already logged in", slog.String("userID", userID))
				session.RedirectUserHome(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
