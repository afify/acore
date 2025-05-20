package db

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"acore/database/pg"

	"github.com/jackc/pgx/v5"
)

type CallFuncParams struct {
	FuncName string
	FuncArgs []interface{}
}

func buildQuery(fnName string, args []interface{}) (string, []interface{}) {
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	sql := fmt.Sprintf("SELECT * FROM %s(%s)",
		fnName,
		strings.Join(placeholders, ","),
	)

	slog.Info("buildQuery", slog.String("query", sql), slog.Any("args", args))
	return sql, args
}

func CallFuncSingle[T any](cfg CallFuncParams) (*T, error) {
	slog.Info("CallFuncSingle started", slog.String("function", cfg.FuncName))

	sql, finalArgs := buildQuery(cfg.FuncName, cfg.FuncArgs)

	slog.Info("Executing query", slog.String("query", sql), slog.Any("args", finalArgs))
	rows, err := pg.DB.Query(context.Background(), sql, finalArgs...)
	if err != nil {
		slog.Error("Query failed", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		slog.Error("Row collection failed", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("CallFuncSingle completed")
	return &item, nil
}
