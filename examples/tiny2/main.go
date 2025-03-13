package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr-tests/examples/common"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "tiny2-app"
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
	return app
}

// func main() {
// 	ctx := context.Background() // with cancel can be passed thru in your actions
// 	app := prepareApp(
// 		// [!code highlight:2]
// 		// use an option store explicitly, or a dummy store by default
// 		cmdr.WithStore(store.New()),
// 	)
// 	if err := app.Run(ctx); err != nil {
// 		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
// 		os.Exit(app.SuggestRetCode())
// 	} else if rc := app.SuggestRetCode(); rc != 0 {
// 		os.Exit(rc)
// 	}
// }
//
// func prepareApp(opts ...cli.Opt) (app cli.App) {
// 	app = cmdr.New(opts...).
// 		Info("tiny2-app", "0.3.1").
// 		Author("The Example Authors") // .Description(``).Header(``).Footer(``)
//
// 	app.Cmd("jump").
// 		Description("jump command").
// 		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
// 		With(func(b cli.CommandBuilder) {
// 			b.Cmd("to").
// 				Description("to command").
// 				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
// 					// cmd.Set() == cmdr.Store(), cmd.Store() == cmdr.Store()
// 					cs := cmdr.Store().WithPrefix("jump.to")
// 					if cs.MustBool("full") {
// 						println()
// 						println(cmd.Set().Dump())
// 					}
// 					cs2 := cmd.Store()
// 					if cs2.MustBool("full") != cs.MustBool("full") {
// 						logz.Panic("a bug found")
// 					}
// 					println("`jump to` has been invoked, and will return with code '1'.")
// 					app.SetSuggestRetCode(1) // ret code must be in 0-255
// 					return                   // handling command action here
// 				}).
// 				With(func(b cli.CommandBuilder) {
// 					b.Flg("full", "f").
// 						Default(false).
// 						Description("full command").
// 						Build()
// 				})
// 		})
// 	return
// }
