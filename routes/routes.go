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

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("views/static")),
		),
	)

	// Public
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			landing.NotFound(w, r)
			return
		}
		landing.Index(w, r)
	})

	mux.HandleFunc("/ping", landing.HeartBeat)
	mux.Handle("/login", mw.PublicOnly(http.HandlerFunc(auth.Login)))
	mux.Handle("/signup", mw.PublicOnly(http.HandlerFunc(auth.Signup)))
	mux.Handle("/auth/google", mw.PublicOnly(http.HandlerFunc(auth.GoogleLogin)))
	mux.Handle("/auth/google/callback", mw.PublicOnly(http.HandlerFunc(auth.GoogleCallback)))

	// Private
	mux.Handle("/home", mw.AuthRequired(http.HandlerFunc(user.Dashboard)))

	return mux
}
