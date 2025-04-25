package main

import (
	"context"
	"io"
	"os"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr-loaders/local"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/dir"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
)

func main() {
	ctx := context.Background() // with cancel can be passed thru in your actions
	app := prepareApp(
		// use an option store explicitly, or a dummy store by default
		cmdr.WithStore(store.New()),

		// [!code highlight:8]
		// import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
		cmdr.WithExternalLoaders(
			local.NewConfigFileLoader(
				local.WithAlternateWriteBack(false), // disable write-back the modified state into alternative config file
				local.WithAlternateDotPrefix(false), // use '<app-name>.toml' instead of '.<app-name>.toml'
			),
			local.NewEnvVarLoader(),
		),

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
			logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}),

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		// The envvars FORCE_DEFAULT_ACTION & FORCE_RUN can override this.
		// cmdr.WithForceDefaultAction(true),

		// cmdr.WithAutoEnvBindings(true),
	)
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("tiny3-app", "0.3.1").
		Author("The Example Authors") // .Description(``).Header(``).Footer(``)

	// another way to disable `cmdr.WithForceDefaultAction(true)` is using
	// env-var FORCE_RUN=1 (builtin already).
	app.Flg("no-default").
		Description("disable force default action").
		// Group(cli.UnsortedGroup).
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				// disable/enable the final state about 'force default action'
				f.Set().Set("app.force-default-action", b)
			}
			return
		}).
		Build()

	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
		Deprecated(`v1.1.0`).
		// Group(cli.UnsortedGroup).
		Hidden(false).
		// Both With(cb) and Build() to end a building sequence
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					// cmd.Set() == cmdr.Set(), cmd.Store() == cmdr.Store()
					cmd.Set().Set("tiny3.working", dir.GetCurrentDir())
					println()
					println(cmd.Set().WithPrefix("tiny3").MustString("working"))
					println(cmdr.Set("tiny3").MustString("working"))

					cs := cmdr.Store().WithPrefix("jump.to")
					if cs.MustBool("full") {
						println()
						println(cmd.Set().Dump())
					}
					cs2 := cmd.Store()
					if cs2.MustBool("full") != cs.MustBool("full") {
						logz.Panic("a bug found")
					}
					app.SetSuggestRetCode(1) // ret code must be in 0-255
					return                   // handling command action here
				}).
				With(func(b cli.CommandBuilder) {
					b.Flg("full", "f").
						Default(false).
						Description("full command").
						Build()
				})
		})

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	app.Cmd("wrong").
		Description("a wrong command to return error for testing").
		// cmdline `FORCE_RUN=1 go run ./tiny wrong -d 8s` to verify this command to see the returned application error.
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			dur := cmd.Store().MustDuration("duration")
			println("the duration is:", dur.String())

			ec := errors.New()
			defer ec.Defer(&err) // store the collected errors in native err and return it
			ec.Attach(io.ErrClosedPipe, errors.New("something's wrong"), os.ErrPermission)
			// see the application error by running `go run ./tiny/tiny/main.go wrong`.
			return
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("duration", "d").
				Default("5s").
				Description("a duration var").
				Build()
		})
	return
}
