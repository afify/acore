package user

import (
	"net/http"

	"acore/models/session"
	modelUser "acore/models/user"

	"acore/render"
)

type DashboardData struct {
	Title string
	User  modelUser.User
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		session.RedirectLogin(w, r)
		return
	}

	title := "Dashboard"
	temp := "user_dashboard.html"

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

	data := DashboardData{
		Title: title,
		User:  *u,
	}

	render.ShowPage(w,
		render.Page[DashboardData]{
			Title:    title,
			PageData: data,
		},
		temp, http.StatusOK,
	)
}
