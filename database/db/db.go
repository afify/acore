package db

import (
	"acore/database/pg"
	"context"
	"fmt"
	"strings"
	"time"
)

func CallFunc(out interface{}, fnName string, args ...interface{}) error {
	fmt.Printf("\x1b[34mDB CALL → %s(%d args)\x1b[0m\n", fnName, len(args))

	ph := make([]string, len(args))
	for i := range args {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT %s(%s)", fnName, strings.Join(ph, ","))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := pg.DB.QueryRowContext(ctx, query, args...)
	return row.Scan(out)
}
