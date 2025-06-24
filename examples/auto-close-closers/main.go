package main

import (
	"context"
	"fmt"
	"os"

	loaders "github.com/hedzr/cmdr-loaders/lite"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	"github.com/hedzr/is/basics"
	logz "github.com/hedzr/logg/slog"

	"database/sql"

	_ "github.com/lib/pq"
)

type dbConn struct {
	conn *sql.DB
}

func (s *dbConn) Close() {
	// here's cleanup operations to free the conn object
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			cmdr.Recycle(err)
		} else {
			s.conn = nil
			logz.Info(`database connection closed`)
		}
	}
}

func (s *dbConn) Open(ctx context.Context) (err error) {
	// do stuffs to open connection to database here

	// locate to `app.resources.db.postgres`
	conf := cmdr.Set().WithPrefix("resources.db.postgres")
	host := conf.MustString("host", "127.0.0.1")
	port := conf.MustInt("port", 5432)
	user := conf.MustString("user", "postgres")
	password := conf.MustString("password", "postgres")
	dbname := conf.MustString("db", "postgres")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	logz.Info(`opening database...`, "connString", psqlInfo)

	s.conn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}

	logz.Info(`database connection opened`)
	err = s.conn.Ping()
	return
}

func (s *dbConn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection not ready")
	}
	return s.conn.QueryContext(ctx, query, args...)
}

func main() {
	// loaders.Create() will load config files if found.
	app := loaders.Create(appName, version, author, desc).WithOpts(
		cmdr.WithPeripherals(map[string]basics.Peripheral{"db": &dbConn{}}),
	).With(func(app cli.App) {
		logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode())
		app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			db := cmdr.PeripheralT[*dbConn]("db")
			if db != nil {
				var rows *sql.Rows
				if rows, err = db.QueryContext(ctx, `SELECT name FROM info`); err == nil {
					defer rows.Close()
					if rows.Next() {
						var str string
						if err = rows.Scan(&str); err == nil {
							println(fmt.Sprintf("info: %s", str))
							return
						}
					}
				}
			}
			println("onAction")
			return
		})
	}).WithAdders(
	// examples.AddTypedFlags,
	).Build()

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

const (
	appName = "auto-close-closers"
	desc    = `a sample to show u how to manage the resource objects`
	version = cmdr.Version
	author  = `The Example Authors`
)
