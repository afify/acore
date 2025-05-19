package session

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"acore/database/db"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Token     string
	ExpiresAt time.Time
}

const (
	CookieName   string = "Acore-Session"
	UserIDHeader string = "X-User-ID"
)

func CreateSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	sess, err := newSession(userID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	setSessionCookie(w, sess.Token, sess.ExpiresAt)
	return nil
}

func newSession(userID uuid.UUID, ipAddress, userAgent string) (*Session, error) {
	token, err := GenerateSessionToken(userID, 24*time.Hour)
	if err != nil {
		slog.Error("newSession", "error", err)
		return nil, fmt.Errorf("newSession(token): %w", err)
	}
	expires := time.Now().Add(24 * time.Hour)

	s := &Session{
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Token:     token,
		ExpiresAt: expires,
	}

	sessionId, err := db.CallFuncSingle[uuid.UUID]("create_user_session", s)
	if err != nil {
		slog.Error("newSession", "error", err)
		return nil, fmt.Errorf("newSession(CallFunc): %w", err)
	}

	s.ID = *sessionId

	slog.Info("newSession", "sessionId", sessionId)
	return s, nil
}

func setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	slog.Info("setting session cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
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
		Name:     CookieName,
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
