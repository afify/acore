package session

import (
	"acore/models/auth"
	"log/slog"
	"net/http"
)

func RedirectLogin(w http.ResponseWriter, r *http.Request) {
	ClearSessionCookie(w, SessionCookieName)
	ClearSessionCookie(w, auth.OauthStateCookieName) // Google oauth
	slog.Info("Redirect to login")
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func RedirectUserHome(w http.ResponseWriter, r *http.Request) {
	slog.Info("Redirect to User Home")
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}
