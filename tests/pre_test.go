package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hedzr/cmdr-loaders/local"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/examples"
	"github.com/hedzr/cmdr/v2/cli/worker"
)

// func TestWorkerS_Pre(t *testing.T) {
// 	app, ww := cleanApp(t, true)

// 	// app := buildDemoApp()
// 	// ww := postBuild(app)
// 	// ww.InitGlobally()
// 	// assert.EqualTrue(t, ww.Ready())
// 	//
// 	// ww.ForceDefaultAction = true
// 	// ww.wrHelpScreen = &discardP{}
// 	// ww.wrDebugScreen = os.Stdout
// 	// ww.wrHelpScreen = os.Stdout

// 	ww.setArgs([]string{"--debug"})
// 	ww.Config.Store = store.New()
// 	ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
// 	_ = app

// 	err := ww.Run(withTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
// 		root.SelfAssert()
// 		t.Logf("root.SelfAssert() passed.")
// 		return
// 	}))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestWorkerS_Pre_v1(t *testing.T) {
	app, ww := cleanApp(t,
		cli.WithExternalLoaders(local.NewConfigFileLoader(), local.NewEnvVarLoader()),
		cli.WithHelpScreenWriter(&discardP{}),
		cli.WithDebugScreenWriter(os.Stdout),
		cli.WithForceDefaultAction(true),
		cli.WithStore(store.New()),
		cli.WithArgs("--debug", "-v"),
	)

	// ww.setArgs([]string{"--debug", "-v"})
	// ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app
	ctx := context.Background()
	err := ww.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestWorkerS_Pre_v3(t *testing.T) {
// 	app, ww := cleanApp(t, true)

// 	ww.setArgs([]string{"--debug", "-vv", "--verbose"})
// 	ww.Config.Store = store.New()
// 	ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
// 	_ = app

// 	err := ww.Run()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

//

//

//

func cleanApp(t *testing.T, opts ...cli.Opt) (app cli.App, ww cli.Runner) { //nolint:revive
	app = buildDemoApp(opts...)
	ww = postBuild(app)
	ctx := context.Background()
	ww.InitGlobally(ctx)
	assertTrue(t, ww.Ready())

	// ww.wrHelpScreen = &discardP{}
	// if helpScreen {
	// 	ww.wrHelpScreen = os.Stdout
	// }
	// ww.wrDebugScreen = os.Stdout
	// ww.ForceDefaultAction = true
	// ww.tasksAfterParse = []taskAfterParse{func(w *workerS, ctx *parseCtx, errParsed error) (err error) { return }}

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

func buildDemoApp(opts ...cli.Opt) (app cli.App) { //nolint:revive
	// cfg := cli.New(cli.WithStore(store.New()))

	cfg := cli.NewConfig(opts...)

	w := worker.New(cfg)

	app = builder.New(w).
		Info("demo-app", "0.3.1").
		Author("hedzr")

	b := app.Cmd("jump").
		Description("jump is a demo command").
		Examples(`jump example`).
		Deprecated(`since: v0.9.1`).
		Hidden(false)
	b1 := b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(``).
		Hidden(false).
		OnAction(func(ctx context.Context, cmd *cli.Command, args []string) (err error) {
			return // handling command action here
		})
	b1.Flg("full", "f").
		Default(false).
		Description("full command").
		Build()
	b1.Build()
	b.Flg("empty", "e").
		Default(false).
		Description("empty command").
		Build()

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

	examples.AttachMoreCommandsForTest(app.Cmd("more", "m"), true)

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
		if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
			r.Root().EnsureTree(context.TODO(), app, r.Root())
			if sra, ok := ww.(interface {
				SetRoot(root *cli.RootCommand, args []string)
			}); ok {
				sra.SetRoot(r.Root(), nil)
			}
		}
	}
	return
}

type discardP struct{}

func (*discardP) Write([]byte) (n int, err error) { return }

func (*discardP) WriteString(string) (n int, err error) { return }

//

//

//

func assertTrue(t testing.TB, cond bool, msg ...any) {
	if cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertTrue failed: %s", mesg)
}

func assertFalse(t testing.TB, cond bool, msg ...any) {
	if !cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertFalse failed: %s", mesg)
}
