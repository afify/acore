package user

import (
	"net/http"

	"acore/models/session"
	modelUser "acore/models/user"

	"acore/render"
)

type UserHomeView struct {
	modelUser.User
	Error string
}

func renderUserHome(w http.ResponseWriter, hv UserHomeView, code int) {
	render.Render(render.RenderRequest{
		Writer:     w,
		Template:   "user_home.html",
		Data:       hv,
		StatusCode: code,
	})
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(session.UserIDHeader)
	if userID == "" {
		session.RedirectLogin(w, r)
		return
	}

	u, err := modelUser.GetByID(userID)
	if err != nil {
		session.RedirectLogin(w, r)
		return
	}

	view := UserHomeView{
		User:  *u,
		Error: "",
	}
	renderUserHome(w, view, http.StatusOK)

}
