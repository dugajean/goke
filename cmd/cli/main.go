package main

import (
	"os"

	app "github.com/dugajean/goke/internal"
)

func main() {
	p := app.Parser{}
	p.Bootstrap()

	l := app.Lockfile{Files: p.FilePaths}
	l.Bootstrap()

	e := app.NewExecuter(p, l)
	if len(os.Args) > 1 {
		e.Execute(os.Args[1], true)
	} else {
		e.Execute("build", true)
	}
}
