package render

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

var templates *template.Template

type RenderRequest struct {
	Writer     http.ResponseWriter
	Template   string
	Data       any
	Headers    http.Header
	StatusCode int
}

type Page[T any] struct {
	Title    string
	PageData T
	Errors   map[string]string
	Warnings map[string]string
}

func InitTemplates() {
	templates = template.Must(template.ParseGlob("views/**/*.html"))
}

func Render(req RenderRequest) {
	w := req.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for k, vs := range req.Headers {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, req.Template, req.Data); err != nil {
		slog.Error("Render: template execution failed",
			"template", req.Template,
			"error", err)
		http.Error(w,
			"Template rendering error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if req.StatusCode != http.StatusOK {
		w.WriteHeader(req.StatusCode)
	}
	if _, err := buf.WriteTo(w); err != nil {
		slog.Error("Render: writing response failed", "error", err)
	}
}

func ShowPage[T any](w http.ResponseWriter, pageData Page[T], tmpl string, code int) {
	Render(RenderRequest{
		Writer:     w,
		Template:   tmpl,
		Data:       pageData,
		StatusCode: code,
	})
}
