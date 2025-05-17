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

	mux.HandleFunc("GET /", landing.Home)
	mux.HandleFunc("GET /ping", landing.HeartBeat)

	// Public Only
	mux.Handle("POST /login", mw.PublicOnly(http.HandlerFunc(auth.SignInForm)))
	mux.Handle("GET /login", mw.PublicOnly(http.HandlerFunc(auth.SignInPage)))
	mux.Handle("POST /signup", mw.PublicOnly(http.HandlerFunc(auth.SignUpForm)))
	mux.Handle("GET /signup", mw.PublicOnly(http.HandlerFunc(auth.SignUpPage)))

	// Private Only
	mux.Handle("GET /home", mw.AuthRequired(http.HandlerFunc(user.UserHome)))
	//mux.Handle("GET /profile", mw.PublicOnly(http.HandlerFunc(user.UserProfile)))

	return mux
}
