package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"acore/database/pg"

	"github.com/jackc/pgx/v5"
)

func buildQuery(fnName string, param interface{}) (sql string, args []interface{}) {
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath == "" { // exported
				args = append(args, v.Field(i).Interface())
			}
		}
	} else {
		args = []interface{}{param}
	}

	ph := make([]string, len(args))
	for i := range ph {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf("SELECT * FROM %s(%s)", fnName, strings.Join(ph, ",")), args
}

// func CallFunc[T any](fnName string, param interface{}, mode CallMode) (interface{}, error) {
// 	switch mode {
// 	case ModeSingle:
// 		return CallFuncSingle[T](fnName, param)
// 	case ModeMulti:
// 		return CallFuncMulti[T](fnName, param)
// 	default:
// 		return nil, fmt.Errorf("invalid CallMode %d", mode)
// 	}
// }

func CallFuncMulti[T any](fnName string, param interface{}) ([]T, error) {
	sql, args := buildQuery(fnName, param)

	rows, err := pg.DB.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// CallFuncSingle invokes the Postgres function fnName(param) and returns exactly one *T.
func CallFuncSingle[T any](fnName string, param interface{}) (*T, error) {
	sql, args := buildQuery(fnName, param)

	rows, err := pg.DB.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}
	return &item, nil
}
