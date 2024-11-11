package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr-loaders/local"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/is"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/examples"
)

const appName = "large-app"
const version = "v1.2.5"
const Version = version

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := prepareApp(
		cmdr.WithStore(store.New()),
		cmdr.WithExternalLoaders(
			local.NewConfigFileLoader(),
			local.NewEnvVarLoader(),
		),

		// cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
		// 	logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
		// 	return
		// }), // cmdr.WithTasksBeforeParse(), cmdr.WithTasksBeforeRun(), cmdr.WithTasksAfterRun

		// // true for debug in developing time, it'll disable onAction on each Cmd.
		// // for productive mode, comment this line.
		// // The envvars FORCE_DEFAULT_ACTION & FORCE_RUN can override this.
		// cmdr.WithForceDefaultAction(true),

		// cmdr.WithSortInHelpScreen(true),       // default it's false
		// cmdr.WithDontGroupInHelpScreen(false), // default it's false
		//
		// cmdr.WithAutoEnvBindings(true),
	)

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	}
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	// A cmdr app will close all peripherals in basics.Closers() at exiting.
	// So you could always register the objects which wanna be cleanup at
	// app terminating, by [basics.RegisterPeripheral(...)].
	// See also: https://github.com/hedzr/is/blob/master/basics/ and
	// is.Closers(), basics.Closers(), ....
	app = cmdr.New(opts...).
		Info(appName, version).
		Author("The Example Authors") // .Description(``).Header(``).Footer(``)

	// another way to disable `cmdr.WithForceDefaultAction(true)` is using
	// env-var FORCE_RUN=1 (builtin already).
	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				f.Set().Set("app.force-default-action", b) // disable/enable the final state about 'force default action'
			}
			return
		}).
		Build()

	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
		Deprecated(`v1.1.0`).
		// Group(cli.UnsortedGroup).
		Hidden(false, false).
		OnEvaluateSubCommands(onEvalJumpSubCommands).
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				// Group(cli.UnsortedGroup).
				Hidden(false).
				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					// cmd.Set() == cmdr.Store():   the whole store
					// cmd.Store() == cmdr.Store(): the sub-store at 'app.cmd', this child-tree holds the live data of command-line args (flags)
					cmd.Set().Set("app.demo.working", dir.GetCurrentDir())
					println()
					println(cmd.Set().WithPrefix("app.demo").MustString("working"))

					cs := cmdr.Store().WithPrefix("jump.to")
					if cs.MustBool("full") {
						println()
						println(cmd.Set().Dump())
					}
					cs2 := cmd.Store()
					if cs2.MustBool("full") != cs.MustBool("full") {
						logz.Panic("a bug found") // cs & cs2 shall point to same a trie-tree-node.
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
		Group(cli.UnsortedGroup).
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
			// see the application error by running `go run ./examples/large/ wrong`.
			return
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("duration", "d").
				Default("5s").
				Description("a duration var").
				Build()
		})

	app.Cmd("consul", "c").
		Description("command set for consul operations").
		With(func(b cli.CommandBuilder) {
			b.Flg("data-center", "dc", "datacenter").
				// Description("set data-center").
				Default("dc-1").
				Build()
		})

	examples.AttachServerCommand(app.Cmd("server", "s"))

	examples.AttachKvCommand(app.Cmd("kv", "kv"))

	examples.AttachMsCommand(app.Cmd("ms", "ms"))

	examples.AttachMoreCommandsForTest(app.Cmd("more", "m"), false)

	app.Cmd("display", "da").
		Description("command set for display adapter operations").
		With(func(b cli.CommandBuilder) {

			b.Cmd("voodoo", "vd").
				Description("command set for voodoo operations").
				With(func(b cli.CommandBuilder) {
					b.Flg("data-center", "dc", "datacenter").
						Default("dc-1").
						Build()
				})

			b.Cmd("nvidia", "nv").
				Description("command set for nvidia operations").
				With(func(b cli.CommandBuilder) {
					b.Flg("data-center", "dc", "datacenter").
						Default("dc-1").
						Build()
				})

			b.Cmd("amd", "amd").
				Description("command set for AMD operations").
				With(func(b cli.CommandBuilder) {
					b.Flg("data-center", "dc", "datacenter").
						Default("dc-1").
						Build()
				})
		})

	return
}

var onceDev sync.Once
var devMode bool

func init() {
	// onceDev is a redundant operation, but we still keep it to
	// fit for defensive programming style.
	onceDev.Do(func() {
		log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC | log.Lshortfile | log.Lmicroseconds)
		log.SetPrefix("")

		if dir.FileExists(".dev-mode") {
			devMode = true
		} else if dir.FileExists("go.mod") {
			data, err := os.ReadFile("go.mod")
			if err != nil {
				return
			}
			content := string(data)

			// dev := true
			if strings.Contains(content, "github.com/hedzr/cmdr/v2/pkg/") {
				devMode = false
			}

			// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
			if strings.Contains(content, "module github.com/hedzr/cmdr") {
				devMode = true
			}
		}

		if devMode {
			is.SetDebugMode(true)
			logz.SetLevel(logz.DebugLevel)
			logz.Debug(".dev-mode file detected, entering Debug Mode...")
		}

		if is.DebugBuild() {
			is.SetDebugMode(true)
			logz.SetLevel(logz.DebugLevel)
		}

		if is.VerboseBuild() {
			is.SetVerboseMode(true)
			if logz.GetLevel() < logz.InfoLevel {
				logz.SetLevel(logz.InfoLevel)
			}
			if logz.GetLevel() < logz.TraceLevel {
				logz.SetLevel(logz.TraceLevel)
			}
		}
	})
}
