package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hedzr/cmdr-tests/examples/common"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "lots-subcmds"
	desc    = ``
)

func main() {
	app := chain(common.PrepareApp(
		appName, desc,
	)(cmd.Commands...))

	ctx := context.Background() // with cancel can be passed thru in your actions
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func chain(app cli.App) cli.App {
	app.Cmd("subcmd").
		Description("subcommands here").
		With(func(b cli.CommandBuilder) {
			for i := 0; i < 5; i++ {
				b.Cmd(fmt.Sprintf("sub-%d", i)).
					Description(fmt.Sprintf("subcommand-%d here", i)).
					With(func(b cli.CommandBuilder) {
						for j := 0; j < 3; j++ {
							b.Cmd(fmt.Sprintf("sub-%d-sub-%d", i, j)).
								Description(fmt.Sprintf("sub-%d-subcommand=%d here", i, j)).
								With(func(b cli.CommandBuilder) {})
						}
					})
			}
		})
	return app
}
