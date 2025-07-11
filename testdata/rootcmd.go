package testdata

import (
	"context"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func BuildCommands(helpScreen bool) *cli.RootCommand {
	app, _ := cleanApp(helpScreen)
	return app.RootCommand()
}

func BuildApp(helpScreen bool) (app cli.App) {
	app, _ = cleanApp(helpScreen)
	return
}

func cleanApp(helpScreen bool) (app cli.App, ww cli.Runner) {
	ctx := context.Background()
	app = buildDemoApp()
	ww = postBuild(app)
	ww.InitGlobally(ctx)

	// assert.EqualTrue(t, ww.Ready())

	// ww.wrHelpScreen = &discardP{}
	// if helpScreen {
	// 	ww.wrHelpScreen = os.Stdout
	// }
	// ww.wrDebugScreen = os.Stdout
	// ww.ForceDefaultAction = true
	// ww.tasksAfterParse = []taskAfterParse{func(w *workerS, ctx *parseCtx) (err error) { return }}

	// ww.setArgs([]string{"--debug"})
	// err := ww.Run(withTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
	// 	root.SelfAssert()
	// 	t.Logf("root.SelfAssert() passed.")
	// 	return
	// }))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ww.TasksBeforeParse = nil
	return
}

func buildDemoApp() (app cli.App) {
	// cfg := cli.New(cli.WithStore(store.New()))

	cfg := cli.NewConfig(
		cli.WithForceDefaultAction(true),
		cli.WithTasksBeforeParse(),
		cli.WithTasksBeforeRun(),
		cli.WithStore(nil),
	)

	w := worker.New(cfg).With(
		worker.WithHelpScreenSets(false, true),
		worker.WithConfig(cfg),
	)

	app = builder.New(w).
		Info("demo-app", "0.3.1").
		Author("hedzr")

	b := app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`).
		Deprecated(`jump is a demo command`).
		Hidden(false)
	b1 := b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(``).
		Hidden(false).
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
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

	common.AttachServerCommand(app.Cmd("server"))

	// app.AddCmd(func(b cli.CommandBuilder) {
	// 	common.AttachKvCommand(b)
	// })

	common.AttachMsCommand(app.Cmd("ms"))

	// app.AddCmd(func(b cli.CommandBuilder) {
	// 	common.AttachMoreCommandsForTest(b)
	// })

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

func postBuild(app cli.App) (ww cli.Runner) {
	if sr, ok := app.(interface{ Worker() cli.Runner }); ok {
		ww = sr.Worker()
		// if ww, ok = sr.Worker().(*workerS); ok {
		if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
			ctx := context.Background()
			r.Root().EnsureTree(ctx, app, r.Root())
			if sr2, ok := ww.(interface {
				SetRoot(root *cli.RootCommand, args []string)
			}); ok {
				sr2.SetRoot(r.Root(), nil)
			}
		}
		// }
	}
	return
}
