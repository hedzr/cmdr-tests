package main

import (
	"context"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
)

func main() {
	app := cmdr.New().
		Info("tiny0-app", "0.3.1").
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

	ctx := context.Background() // with cancel can be passed thru in your actions
	if err := app.Run(ctx); err != nil {
		println("Application Error:", err)
	}
}
