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
	appName = "lots-flags"
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
	app.Cmd("cmd").
		Description("subcommands with lots of flags here").
		With(func(b cli.CommandBuilder) {
			for i := 0; i < 15; i++ {
				b.Flg(fmt.Sprintf("option-%d", i), fmt.Sprintf("o%d", i)).
					Description(fmt.Sprintf("flag(option)-%d here", i)).
					Build()
			}
		})
	return app
}
