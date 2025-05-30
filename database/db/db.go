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
	FuncArgs []any
}

func buildQuery(fnName string, args []any) (string, []any) {
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	sql := fmt.Sprintf("SELECT * FROM %s(%s)", fnName, strings.Join(placeholders, ","))

	slog.Info("buildQuery", slog.String("query", sql), slog.Any("args", args))
	return sql, args
}

func CallFuncSingle[T any](cfg CallFuncParams) (*T, error) {
	sql, args := buildQuery(cfg.FuncName, cfg.FuncArgs)

	row, err := pg.DB.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	item, err := pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}
	return &item, nil
}
