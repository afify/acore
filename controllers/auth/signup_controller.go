package auth

import (
	"acore/models/auth"
	"acore/models/session"
	"acore/models/user"
	"acore/models/validator"
	"acore/render"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func renderSignupPage(w http.ResponseWriter, form auth.SignupReq, errs map[string]string, status int) {
	render.ShowPage(w, render.Page[auth.SignupReq]{
		Title:    "Sign Up",
		PageData: form,
		Errors:   errs,
	}, "signup.html", status)
}

func handleSignupPost(w http.ResponseWriter, r *http.Request) {
	form, ok := parseSignupRequest(w, r)
	if !ok {
		return
	}

	user, errs := createUserFromForm(form)
	if errs != nil {
		renderSignupPage(w, form, errs, http.StatusConflict)
		return
	}

	if err := session.CreateSession(w, r, user.ID, session.SessionTypeWeb, auth.AuthProviderEmail); err != nil {
		renderSignupPage(w, form,
			map[string]string{"Signup": "Could not create session"},
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func createUserFromForm(form auth.SignupReq) (*user.User, map[string]string) {
	hash, err := auth.HashPassword(form.Password)
	if err != nil {
		slog.Error("HashPassword failed", "error", err)
		return nil, map[string]string{"Signup": "Could not create user"}
	}
	form.Password = hash

	user, err := auth.CreateUser(form)
	if err != nil {
		return nil, mapPgError(err)
	}
	if user == nil {
		return nil, map[string]string{"Signup": "Could not create user"}
	}
	return user, nil
}

func mapPgError(err error) map[string]string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return map[string]string{"Email": "That email is already registered"}
		case "users_username_key":
			return map[string]string{"Username": "That username is already in use"}
		}
	}
	slog.Error("CreateUser failed", "error", err)
	return map[string]string{"Signup": "Could not create user"}
}

func parseSignupRequest(w http.ResponseWriter, r *http.Request) (auth.SignupReq, bool) {
	form, errs, err := validator.BindAndValidateForm[auth.SignupReq](r)
	if err != nil {
		renderSignupPage(w, form,
			map[string]string{"Signup": "Invalid request"},
			http.StatusBadRequest,
		)
		return form, false
	}
	if len(errs) > 0 {
		renderSignupPage(w, form, errs, http.StatusUnprocessableEntity)
		return form, false
	}
	return form, true
}
