package middleware

import (
	"log/slog"
	"net/http"

	"acore/models/session"
)

func AuthRequired(next http.Handler) http.Handler {
	slog.Info("AuthRequired")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("AuthRequired: incoming request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", r.Header.Get("X-Forwarded-For")),
			slog.String("user-agent", r.UserAgent()),
		)

		c, err := r.Cookie(session.SessionCookieName)
		if err != nil || c.Value == "" {
			slog.Error("Get Cookie:", "error", err)
			session.RedirectLogin(w, r)
			return
		}

		userID, sessionID, err := session.ValidateSessionToken(c.Value)
		if err != nil {
			slog.Error("Verify Session:", "error", err)
			session.RedirectLogin(w, r)
			return
		}

		session.SetSessionCookieStrict(w,
			session.SessionCookieName,
			c.Value,
			session.DefaultExpiry(),
		)

		r.Header.Set(session.UserIDHeader, userID.String())
		r.Header.Set(session.SessionIDHeader, sessionID.String())
		slog.Info("Session valid, setting user ID header",
			"userID", userID,
			"sessionID", sessionID,
		)
		next.ServeHTTP(w, r)
	})
}

func PublicOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("PublicOnly: incoming request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", r.Header.Get("X-Forwarded-For")),
			slog.String("user-agent", r.UserAgent()),
		)
		if c, err := r.Cookie(session.SessionCookieName); err == nil && c.Value != "" {
			if userID, sessionID, err := session.ValidateSessionToken(c.Value); err == nil {
				slog.Info("Session valid, setting user ID header",
					"userID", userID,
					"sessionID", sessionID,
				)
				session.RedirectUserHome(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
