package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr-loaders/local"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/examples"
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
		Info("large-app", "0.3.1").
		Author("hedzr")

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

	b = app.Cmd("consul", "c").
		Description("command set for consul operations")
	b.Flg("data-center", "dc", "datacenter").
		// Description("set data-center").
		Default("dc-1").
		Build()
	b.Build()

	examples.AttachServerCommand(app.Cmd("server", "s"))

	examples.AttachKvCommand(app.Cmd("kv", "kv"))

	examples.AttachMsCommand(app.Cmd("ms", "ms"))

	examples.AttachMoreCommandsForTest(app.Cmd("more", "m"), false)

	b = app.Cmd("display", "da").
		Description("command set for display adapter operations")

	b1 = b.Cmd("voodoo", "vd").
		Description("command set for voodoo operations")
	b1.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b1.Build()

	b2 := b.Cmd("nvidia", "nv").
		Description("command set for nvidia operations")
	b2.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b2.Build()

	b3 := b.Cmd("amd", "amd").
		Description("command set for AMD operations")
	b3.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b3.Build()

	b.Build()

	return
}
