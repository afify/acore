package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type responseLogger struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.status = code
	rl.ResponseWriter.WriteHeader(code)
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	if rl.status == 0 {
		rl.status = http.StatusOK
	}
	rl.body.Write(b)
	return rl.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read and restore request body
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Capture response
		rl := &responseLogger{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(rl, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rl.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_body", string(reqBody),
			"response_body", rl.body.String(),
		)
	})
}
