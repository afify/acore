package landing

import (
	"fmt"
	"net/http"
	"os"

	"acore/render"
)

type PageData struct {
}

func HeartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong from "+os.Getenv("APP_NAME")+"["+os.Getenv("COMMIT")+"]\n")
}

type LandingView struct {
	Error     string
	Container string
	Commit    string
}

func showLanding(w http.ResponseWriter, v LandingView, code int) {
	render.Render(render.RenderRequest{
		Writer:     w,
		Template:   "index.html",
		Data:       v,
		StatusCode: code,
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	container := os.Getenv("APP_NAME")
	commit := os.Getenv("COMMIT")
	showLanding(w, LandingView{
		Error:     "",
		Container: container,
		Commit:    commit,
	}, http.StatusOK)

}
