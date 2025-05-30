package session

import (
	"acore/models/auth"
	"log/slog"
	"net/http"
	"time"
)

func SetSessionCookieStrict(w http.ResponseWriter, cookieName, token string, expires time.Time) {
	slog.Info("setting session cookie",
		"name", cookieName,
		"token", token,
		"expires", expires,
		"type", http.SameSiteStrictMode,
	)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expires,
	})
}

func SetSessionCookie(w http.ResponseWriter, cookieName, token string, expires time.Time) {
	slog.Info("setting session cookie",
		"name", cookieName,
		"token", token,
		"expires", expires,
		"type", http.SameSiteNoneMode,
	)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expires,
	})
}

func ClearSessionCookie(w http.ResponseWriter, cookieName string) {
	slog.Info("Clear session cookie", "name", cookieName)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
}

func ClearSessionCookies(w http.ResponseWriter) {
	ClearSessionCookie(w, SessionCookieName)
	ClearSessionCookie(w, auth.OauthStateCookieName)
}
