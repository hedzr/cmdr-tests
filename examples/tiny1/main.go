package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
)

func main() {
	ctx := context.Background() // with cancel can be passed thru in your actions
	app := prepareApp()
	if err := app.Run(ctx); err != nil {
		println("Application Error:", err)
		os.Exit(app.SuggestRetCode())
	}
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("tiny1-app", "0.3.1").
		Author("The Example Authors") // .Description(``).Header(``).Footer(``)

	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			println("jump command:", cmd)
			if cmd.FlagBy("full").GetTriggeredTimes() > 0 {
				// for dummy store, `if cmd.Store().MustBool("full") {}` cannot work
				println("Dump", cmd.Set().Dump()) // nothing to display since a dummy store created
			}
			return
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("full", "f").
				Default(false).
				Description("full option here").
				Build()
		})
	return
}
