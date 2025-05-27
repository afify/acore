package session

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"acore/models/auth"

	"github.com/google/uuid"
)

type SessionType int16

const (
	SessionTypeWeb        SessionType = 1
	SessionTypeMobile     SessionType = 2
	SessionTypeAPI        SessionType = 3
	SessionCookieName     string      = "Acore-Session"
	UserIDHeader          string      = "X-User-ID"
	SessionIDHeader       string      = "X-Session-ID"
	SessionExpiryDuration             = 10 * 24 * time.Hour // 10 Days
)

func (st SessionType) String() string {
	switch st {
	case SessionTypeWeb:
		return "web"
	case SessionTypeMobile:
		return "mobile"
	case SessionTypeAPI:
		return "api"
	default:
		return "unknown"
	}
}

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	SessionType  int16     `db:"session_type_id"`
	AuthProvider int16     `db:"auth_provider_id"`
	Token        string    `db:"session_token"`
	IPAddress    string    `db:"ip_address"`
	UserAgent    string    `db:"user_agent"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type SessionParams struct {
	UserID    uuid.UUID
	Type      SessionType
	Provider  auth.AuthProvider
	Token     string
	IPAddress string
	UserAgent string
	ExpiresAt time.Time
}

func getSessionKey() ([]byte, error) {
	key := os.Getenv("SESSION_ENC_KEY")
	if len(key) != 32 {
		return nil, errors.New("SESSION_ENC_KEY must be 32 bytes")
	}
	return []byte(key), nil
}

func DefaultExpiry() time.Time {
	return time.Now().UTC().Add(SessionExpiryDuration)
}

func CreateSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID, st SessionType, provider auth.AuthProvider) error {
	tok, err := GenerateSessionToken(userID)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	params := SessionParams{
		UserID:    userID,
		Type:      st,
		Provider:  provider,
		Token:     tok,
		IPAddress: r.Header.Get("X-Forwarded-For"),
		UserAgent: r.UserAgent(),
		ExpiresAt: DefaultExpiry(),
	}

	sess, err := dbCreateUserSession(params)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}

	// Oauth require SameSite = None
	if provider == auth.AuthProviderEmail {
		SetSessionCookieStrict(w, SessionCookieName, sess.Token, sess.ExpiresAt)
	} else {
		SetSessionCookieOauth(w, SessionCookieName, sess.Token, sess.ExpiresAt)
	}

	slog.Info("NewSession", "session", sess)
	return nil
}
