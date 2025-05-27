package landing

import (
	"fmt"
	"net/http"
	"os"

	"acore/render"
)

func HeartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	commit := os.Getenv("COMMIT")
	fmt.Fprintf(w, "Pong from %s\n", commit)
}

type LandingView struct {
	Container string
	Commit    string
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	title := "Zeerobot - NotFound"
	temp := "404.html"

	w.WriteHeader(http.StatusNotFound)
	render.ShowPage(w,
		render.Page[LandingView]{
			Title: title,
		},
		temp, http.StatusNotFound,
	)
}

func Index(w http.ResponseWriter, r *http.Request) {
	title := "Zeerobot"
	temp := "index.html"
	var pageData LandingView

	pageData.Container = os.Getenv("APP_NAME")
	pageData.Commit = os.Getenv("COMMIT")

	render.ShowPage(w,
		render.Page[LandingView]{
			Title:    title,
			PageData: pageData,
		},
		temp, http.StatusOK,
	)
}
