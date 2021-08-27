package db

import (
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"

	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/dbEngine/dbEngine/psql"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
)

func GetDB(ctxApis apis.CtxApis) *dbEngine.DB {
	conn := psql.NewConn(func(context.Context, *pgx.Conn) error { return nil },
		nil, printNotice)
	ctx := context.WithValue(ctxApis, "dbURL", "")
	ctx = context.WithValue(ctx, "fillSchema", true)
	db, err := dbEngine.NewDB(ctx, conn)
	if err != nil {
		logs.ErrorLog(err, "")
		return nil
	}

	return db
}

func printNotice(c *pgconn.PgConn, n *pgconn.Notice) {
	if n.Severity == "INFO" {
		logs.StatusLog(n.Message)
	} else if n.Code > "00000" {
		err := (*pgconn.PgError)(n)
		logs.ErrorLog(err, n.Hint, err.SQLState(), err.File, err.Line, err.Routine)
	} else {
		logs.ErrorLog(errors.New(n.Message + n.Severity))
	}
}
