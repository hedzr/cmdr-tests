package common

import (
	loaders "github.com/hedzr/cmdr-loaders"
	"github.com/hedzr/cmdr/v2/cli"
)

func PrepareApp(appName, desc string, opts ...cli.Opt) func(adders ...cli.CmdAdder) (app cli.App) {
	return loaders.PrepareApp(appName, desc, opts...)
}
