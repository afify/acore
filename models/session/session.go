package models

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"acore/database/pg"
)

type Session struct {
	ID           string
	UserID       string
	IPAddress    string
	UserAgent    string
	SessionToken string
	ExpiresAt    time.Time
}

func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {
	sess, err := newSession(r.Context(), userID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	setSessionCookie(w, sess.SessionToken, sess.ExpiresAt)
	return nil
}

func newSession(ctx context.Context, userID, ipAddress, userAgent string) (*Session, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("NewSession(rand): %w", err)
	}
	token := base64.URLEncoding.EncodeToString(b)
	expires := time.Now().Add(24 * time.Hour)

	const fn = `SELECT public.create_user_session($1, $2, $3, $4, $5)`
	row := pg.DB.QueryRowContext(ctx, fn,
		userID,
		token,
		ipAddress,
		userAgent,
		expires,
	)

	s := &Session{
		UserID:       userID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		SessionToken: token,
		ExpiresAt:    expires,
	}
	if err := row.Scan(&s.ID); err != nil {
		return nil, fmt.Errorf("NewSession(scan): %w", err)
	}
	return s, nil
}

func setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expires,
	})
}
