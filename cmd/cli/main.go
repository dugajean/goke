package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
)

func main() {
	argIndex := app.PermutateArgs(os.Args)
	opts := app.GetOptions()

	handleGlobalFlags(&opts)

	cfg, err := app.ReadYamlConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Wrappers
	fs := app.LocalFileSystem{}
	proc := app.ShellProcess{}

	// Main components
	p := app.NewParser(cfg, &opts, &fs)
	p.Bootstrap()

	l := app.NewLockfile(p.GetFilePaths(), &opts, &fs)
	l.Bootstrap()

	e := app.NewExecutor(&p, &l, &opts, &proc)
	e.Start(parseTaskName(argIndex))
}
