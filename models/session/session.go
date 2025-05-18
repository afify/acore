package session

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"acore/database/db"
)

type Session struct {
	ID           string
	UserID       string
	IPAddress    string
	UserAgent    string
	SessionToken string
	ExpiresAt    time.Time
}

const sessionCookieName string = "Session_$"

func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {
	sess, err := newSession(userID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	setSessionCookie(w, sess.SessionToken, sess.ExpiresAt)
	return nil
}

func newSession(userID, ipAddress, userAgent string) (*Session, error) {
	var sessionId string

	token, err := GenerateSessionToken(userID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("newSession(token): %w", err)
	}
	expires := time.Now().Add(24 * time.Hour)

	const fn = "public.create_user_session"
	if err := db.CallFunc(
		&sessionId,
		fn,
		userID,
		token,
		ipAddress,
		userAgent,
		expires,
	); err != nil {
		return nil, fmt.Errorf("newSession(CallFunc): %w", err)
	}

	return &Session{
		ID:           sessionId,
		UserID:       userID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		SessionToken: token,
		ExpiresAt:    expires,
	}, nil
}

func setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	slog.Info("setting session cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expires,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	slog.Info("Clear session cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func RedirectLogin(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	slog.Info("Redirect to login")
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func RedirectUserHome(w http.ResponseWriter, r *http.Request) {
	slog.Info("Redirect to User Home")
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}
