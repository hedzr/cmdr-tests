package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr-loaders/local"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
)

func main() {
	app := prepareApp()

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )
	if err := app.Run(ctx); err != nil {
		logz.Error("Application Error:", "err", err)
	}
}

func prepareApp() (app cli.App) {
	app = cmdr.New().
		Info("demo-app", "0.3.1").
		Author("hedzr")
	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			app.Store().Set("app.force-default-action", false)
			return
		}).
		Build()

	b := app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`).
		Deprecated(`jump is a demo command`).
		Hidden(false)

	b1 := b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(`v0.1.1`).
		Hidden(false).
		OnAction(func(ctx context.Context, cmd *cli.Command, args []string) (err error) {
			app.Store().Set("app.demo.working", dir.GetCurrentDir())
			println()
			println(dir.GetCurrentDir())
			println()
			println(app.Store().Dump())
			return // handling command action here
		})
	b1.Flg("full", "f").
		Default(false).
		Description("full command").
		Build()
	b1.Build()

	b.Build()

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	return
}
