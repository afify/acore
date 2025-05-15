package routes

import (
	"acore/controllers"
	"acore/middleware"
	"log"
	"net/http"
	"os"
	"text/template"
)

type PageData struct {
	Color  string // will be "blue" or "green"
	Commit string // short SHA
}

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	data := PageData{
		Color:  os.Getenv("DEPLOY_COLOR"), // set this in your docker-compose for each service
		Commit: os.Getenv("COMMIT"),       // your short SHA exported from Makefile
	}

	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if err := indexTmpl.Execute(w, data); err != nil {
			log.Printf("template execute error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})

	// Public routes
	mux.HandleFunc("/ping", controllers.HeartBeat)
	mux.HandleFunc("/login", controllers.LoginUser)
	mux.HandleFunc("/users", controllers.HandleUsers)

	mux.Handle("/user/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.HandleUsers)))

	return mux
}
