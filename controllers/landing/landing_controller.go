package landing

import (
	"fmt"
	"net/http"
	"os"

	"acore/render"
)

type PageData struct {
	Color  string // will be "blue" or "green"
	Commit string // short SHA
}

func HeartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong from "+os.Getenv("APP_NAME")+"["+os.Getenv("COMMIT")+"]\n")
}

type LandingView struct {
	Error string
	// TODO: add more fields here as you need them
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
	showLanding(w, LandingView{}, http.StatusOK)
}
