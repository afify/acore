package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func Init() {
	w := os.Stderr
	handler := tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelDebug,
		AddSource:  true,
		TimeFormat: time.DateTime,
	})
	slog.SetDefault(slog.New(handler))
	slog.Info("Logger initialized")
}
