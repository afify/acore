package auth

import (
	"log/slog"
	"net/http"

	auth "acore/models/auth"
	session "acore/models/session"
	"acore/render"
)

type SignUpView struct {
	Form  auth.SignUpReq
	Error string
}

type SignInView struct {
	Form  auth.SignInReq
	Error string
}

func bindSignUp(r *http.Request) (auth.SignUpReq, error) {
	if err := r.ParseForm(); err != nil {
		return auth.SignUpReq{}, err
	}
	req := auth.SignUpReq{
		UserName: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	return req, nil
}

func bindSignIn(r *http.Request) (auth.SignInReq, error) {
	if err := r.ParseForm(); err != nil {
		return auth.SignInReq{}, err
	}
	return auth.SignInReq{
		EmailUsername: r.FormValue("email-username"),
		Password:      r.FormValue("password"),
	}, nil
}

func showSignUp(w http.ResponseWriter, form auth.SignUpReq, msg string, code int) {
	render.Render(render.RenderRequest{
		Writer:     w,
		Template:   "signup.html",
		Data:       SignUpView{Form: form, Error: msg},
		Headers:    nil,
		StatusCode: code,
	})
}

func showSignIn(w http.ResponseWriter, form auth.SignInReq, msg string, code int) {
	render.Render(render.RenderRequest{
		Writer:     w,
		Template:   "signin.html",
		Data:       SignInView{Form: form, Error: msg},
		Headers:    nil,
		StatusCode: code,
	})
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	showSignUp(w, auth.SignUpReq{}, "", http.StatusOK)
}

func SignInPage(w http.ResponseWriter, r *http.Request) {
	showSignIn(w, auth.SignInReq{}, "", http.StatusOK)
}

func SignUpForm(w http.ResponseWriter, r *http.Request) {
	form, err := bindSignUp(r)
	if err != nil {
		showSignUp(w, form, "Invalid form submission", http.StatusBadRequest)
		return
	}

	form.Password, _ = auth.HashPassword(form.Password)
	userID, err := auth.CreateUser(form)
	if err != nil {
		slog.Error("SignUpForm failed", "error", err)
		showSignUp(w, form, err.Error(), http.StatusConflict)
		return
	}

	err = session.CreateSession(w, r, userID)
	if err != nil {
		slog.Error("SignUpForm failed", "error", err)
		showSignUp(w, form, err.Error(), http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SignInForm(w http.ResponseWriter, r *http.Request) {
	form, err := bindSignIn(r)
	if err != nil {
		showSignIn(w, form, "Invalid form submission", http.StatusBadRequest)
		return
	}

	userID, err := auth.Authenticate(form)
	if err != nil {
		showSignIn(w, form, "Wrong email/username or password", http.StatusUnauthorized)
		return
	}

	if err := session.CreateSession(w, r, userID); err != nil {
		showSignIn(w, form, "Could not create session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
