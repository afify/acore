package user

import (
	"net/http"

	modelUser "acore/models/user"

	"acore/render"
)

type UserHomeView struct {
	Data  modelUser.User
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
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	u, err := modelUser.GetByID(userID)
	if err != nil {
		renderUserHome(w,
			UserHomeView{
				Data:  modelUser.User{},
				Error: "Unable to load user",
			},
			http.StatusInternalServerError)
		return
	}

	renderUserHome(w, UserHomeView{Data: *u, Error: ""}, http.StatusOK)
}
