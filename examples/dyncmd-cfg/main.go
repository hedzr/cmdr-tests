package main

import (
	"context"
	"os"

	loaders "github.com/hedzr/cmdr-loaders/lite"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

const (
	appName = "dyncmd-cfg"
	desc    = `dyncmd defined in config file`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	ctx := context.Background()

	app := loaders.Create(appName, version, author, desc).
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithAdders(cmd.Commands...).
		Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
