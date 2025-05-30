package auth

import (
	"net/http"

	"acore/models/auth"
	"acore/models/session"
	"acore/models/validator"
	"acore/render"
	"log/slog"

	"github.com/google/uuid"
)

func renderLoginPage(w http.ResponseWriter, form auth.LoginReq, errs map[string]string, status int) {
	render.ShowPage(w, render.Page[auth.LoginReq]{
		Title:    "Login",
		PageData: form,
		Errors:   errs,
		Warnings: map[string]string{},
	}, "login.html", status)
}

func handleLoginPost(w http.ResponseWriter, r *http.Request) {
	form, ok := parseLoginRequest(w, r)
	if !ok {
		return
	}

	userID, errs := processLogin(form)
	if errs != nil {
		renderLoginPage(w, form, errs, http.StatusUnauthorized)
		return
	}

	if err := session.CreateSession(w, r, userID, session.SessionTypeWeb, auth.AuthProviderEmail); err != nil {
		slog.Error("Session creation failed", "error", err)
		renderLoginPage(w, form,
			map[string]string{"Login": "Could not create session"},
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func parseLoginRequest(w http.ResponseWriter, r *http.Request) (auth.LoginReq, bool) {
	form, errs, err := validator.BindAndValidateForm[auth.LoginReq](r)
	if err != nil {
		renderLoginPage(w, form,
			map[string]string{"Signup": "Invalid request"},
			http.StatusBadRequest,
		)
		return form, false
	}
	if len(errs) > 0 {
		renderLoginPage(w, form, errs, http.StatusUnprocessableEntity)
		return form, false
	}
	return form, true
}

func processLogin(form auth.LoginReq) (uuid.UUID, map[string]string) {
	userID, err := auth.Authenticate(form)
	if err != nil {
		slog.Error("Authenticate failed", "error", err)
		return uuid.Nil, map[string]string{"Login": "Wrong credentials"}
	}
	return userID, nil
}
