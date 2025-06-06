package main

import (
	"context"
	"os"

	loaders "github.com/hedzr/cmdr-loaders/lite"
	"github.com/hedzr/cmdr-tests/examples/demo/cmd"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	logz "github.com/hedzr/logg/slog"
)

func main() {
	ctx := context.Background()
	app := prepareApp(cmd.Commands...) // define your own commands implementations with cmd/*.go
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func prepareApp(commands ...cli.CmdAdder) (app cli.App) {
	return loaders.Create(
		AppNameExample, version, author, desc,
		append([]cli.Opt{
			cmdr.WithAutoEnvBindings(true),  // default it's false
			cmdr.WithSortInHelpScreen(true), // default it's false
			// cmdr.WithDontGroupInHelpScreen(false), // default it's false
			// cmdr.WithForceDefaultAction(false),
		})...,
	).
		// importing devmode package and run its init():
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithBuilders(
			common.AddHeadLikeFlagWithoutCmd, // add a `--line` option, feel free to remove it.
			common.AddToggleGroupFlags,       //
			common.AddTypedFlags,             //
			common.AddKilobytesFlag,          //
			common.AddValidArgsFlag,          //
		).
		WithAdders(commands...).
		Build()
}
