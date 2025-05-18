package db

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"acore/database/pg"
)

func CallFunc(out interface{}, fnName string, args ...interface{}) error {
	start := time.Now()

	slog.Info("DB call start",
		slog.String("fn", fnName),
		slog.Int("args", len(args)),
	)

	// build the SQL
	ph := make([]string, len(args))
	for i := range args {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("SELECT %s(%s)", fnName, strings.Join(ph, ","))

	// execute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := pg.DB.QueryRow(ctx, query, args...).Scan(out)

	slog.Info("DB call end",
		slog.String("fn", fnName),
		slog.Int("args", len(args)),
		slog.Duration("duration", time.Since(start)),
		slog.Any("error", err),
	)

	return err
}
