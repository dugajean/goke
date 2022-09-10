package main

import (
	"os"

	app "github.com/dugajean/goke/internal"
	"github.com/dugajean/goke/internal/cli"
)

func main() {
	opts := cli.GetOptions()

	p := app.NewParser(app.ReadYamlConfig(), opts.ClearCache)
	p.Bootstrap()

	l := app.Lockfile{Files: p.FilePaths}
	l.Bootstrap()

	e := app.NewExecuter(&p, &l, &opts)

	arg := ""
	if len(os.Args) > 0 {
		arg = os.Args[1]
	}

	e.Start(arg)
}
