package auth

import (
	"net/http"

	"acore/models/auth"
	"acore/models/session"
)

func Login(w http.ResponseWriter, r *http.Request) {
	session.ClearSessionCookies(w)
	switch r.Method {
	case http.MethodGet:
		renderLoginPage(w, auth.LoginReq{}, nil, http.StatusOK)
	case http.MethodPost:
		handleLoginPost(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	session.ClearSessionCookies(w)
	switch r.Method {
	case http.MethodGet:
		renderSignupPage(w, auth.SignupReq{}, nil, http.StatusOK)
	case http.MethodPost:
		handleSignupPost(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete_user session
	session.RedirectLogin(w, r)
}
