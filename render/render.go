package render

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type RenderRequest struct {
	Writer     http.ResponseWriter
	Template   string
	Data       interface{}
	Headers    http.Header
	StatusCode int
}

var templates *template.Template

func InitTemplates() {
	const pattern = "views/*/*.html"

	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("render.Init: glob failed", "pattern", pattern, "error", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		slog.Error("render.Init: no template files found", "pattern", pattern)
		os.Exit(1)
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		slog.Error("render.Init: failed to parse templates", "error", err)
		os.Exit(1)
	}

	templates = tmpl
}

func Render(req RenderRequest) {
	w := req.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for key, values := range req.Headers {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(req.StatusCode)
	if err := templates.ExecuteTemplate(w, req.Template, req.Data); err != nil {
		slog.Error("Render: template execution failed",
			"template", req.Template,
			"error", err)
		http.Error(w,
			"Template rendering error: "+err.Error(),
			http.StatusInternalServerError)
	}
}
