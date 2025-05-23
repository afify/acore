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
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))),
	)

	// 2) Public pages
	mux.HandleFunc("/", landing.Home)          // GET /
	mux.HandleFunc("/ping", landing.HeartBeat) // GET /ping

	// 3) Login (GET shows page, POST processes form)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mw.PublicOnly(http.HandlerFunc(auth.SignInPage)).ServeHTTP(w, r)
		case http.MethodPost:
			mw.PublicOnly(http.HandlerFunc(auth.SignInForm)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// 4) Signup (same pattern as Login)
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mw.PublicOnly(http.HandlerFunc(auth.SignUpPage)).ServeHTTP(w, r)
		case http.MethodPost:
			mw.PublicOnly(http.HandlerFunc(auth.SignUpForm)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

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
