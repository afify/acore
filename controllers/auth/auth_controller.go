// in controllers/auth/auth_controller.go
package auth

import (
	"log/slog"
	"net/http"

	"acore/models/auth"
	authModel "acore/models/auth"
	"acore/models/session"
	"acore/models/validator"
	"acore/render"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var title string = "Sign In"
	var temp string = "signin.html"
	var form authModel.SignInReq

	session.ClearSessionCookie(w, session.SessionCookieName)
	session.ClearSessionCookie(w, auth.OauthStateCookieName) // Google oauth

	switch r.Method {
	case http.MethodGet:
		render.ShowPage(w,
			render.Page[authModel.SignInReq]{
				Title:    title,
				PageData: form,
			},
			temp, http.StatusOK,
		)

	case http.MethodPost:
		errs, err := validator.BindAndValidate(r, &form)
		if err != nil {
			slog.Error("Validation failed", "error", err, "errors", errs)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if len(errs) > 0 {
			slog.Error("Validation failed", "errors", errs)
			render.ShowPage(w,
				render.Page[authModel.SignInReq]{
					Title:    title,
					PageData: form,
					Errors:   errs,
					Warnings: map[string]string{},
				},
				temp, http.StatusUnprocessableEntity,
			)
			return
		}

		userID, err := authModel.Authenticate(form)
		if err != nil {
			slog.Error("Authenticate failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignInReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Login": "Wrong credentials",
					},
				},
				temp,
				http.StatusUnauthorized,
			)
			return
		}

		if err := session.CreateSession(w, r, userID, session.SessionTypeWeb, authModel.AuthProviderEmail); err != nil {
			slog.Error("Session failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignInReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Login": "Could not create session",
					},
				},
				temp, http.StatusInternalServerError,
			)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var title string = "Sign Up"
	var temp string = "signup.html"
	var form authModel.SignUpReq

	session.ClearSessionCookie(w, session.SessionCookieName)
	session.ClearSessionCookie(w, auth.OauthStateCookieName) // Google oauth
	switch r.Method {
	case http.MethodGet:
		render.ShowPage(w,
			render.Page[authModel.SignUpReq]{
				Title:    title,
				PageData: form,
			},
			temp, http.StatusOK,
		)

	case http.MethodPost:
		errs, err := validator.BindAndValidate(r, &form)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if len(errs) > 0 {
			render.ShowPage(w,
				render.Page[authModel.SignUpReq]{
					Title:    title,
					PageData: form,
					Errors:   errs,
				},
				temp, http.StatusUnprocessableEntity,
			)
			return
		}

		form.Password, err = authModel.HashPassword(form.Password)
		if err != nil {
			slog.Error("HashPassword failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignUpReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Signup": "Could not create user",
					},
				},
				temp, http.StatusConflict,
			)
			return
		}

		u, err := authModel.CreateUser(form)
		if err != nil {
			slog.Error("CreateUser failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignUpReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Signup": "Could not create user",
					},
				},
				temp, http.StatusConflict,
			)
			return
		}

		if u == nil {
			slog.Error("CreateUser failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignUpReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Signup": "Could not create user",
					},
				},
				temp, http.StatusConflict,
			)
			return
		}

		if err := session.CreateSession(w, r, u.ID, session.SessionTypeWeb, authModel.AuthProviderEmail); err != nil {
			slog.Error("Session failed", "error", err)
			render.ShowPage(w,
				render.Page[authModel.SignUpReq]{
					Title:    title,
					PageData: form,
					Errors: map[string]string{
						"Signup": "Could not create user",
					},
				},
				temp, http.StatusConflict,
			)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete_user session
	session.RedirectLogin(w, r)
}
