package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
	"github.com/dugajean/goke/internal/cli"
)

func main() {
	argIndex := app.PermutateArgs(os.Args)
	opts := cli.GetOptions()

	handleGlobalFlags(&opts)

	cfg, err := app.ReadYamlConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fs := app.LocalFileSystem{}
	p := app.NewParser(cfg, &opts, &fs)
	p.Bootstrap()

	l := app.NewLockfile(p.FilePaths, &opts, &fs)
	l.Bootstrap()

	e := app.NewExecutor(&p, &l, &opts)
	e.Start(parseTaskName(argIndex))
}
