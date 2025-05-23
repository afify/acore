package routes

import (
	"net/http"

	"acore/controllers/auth"
	"acore/controllers/landing"
	"acore/controllers/user"
	mw "acore/middleware"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 1) Static assets (CSS, JS, images)
	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))),
	)

	mux.HandleFunc("/", landing.Home)
	mux.HandleFunc("/ping", landing.HeartBeat)
	mux.Handle("/login", mw.PublicOnly(http.HandlerFunc(auth.Login)))
	mux.Handle("/signup", mw.PublicOnly(http.HandlerFunc(auth.Signup)))

	// 5) Protected user home (only GET)
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		mw.AuthRequired(http.HandlerFunc(user.UserHome)).ServeHTTP(w, r)
	})

	return mux
}
