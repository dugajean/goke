package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
)

func parseTaskName(argIndex int) string {
	arg := ""

	if len(os.Args) > argIndex {
		arg = os.Args[argIndex]
	}

	return arg
}

func handleGlobalFlags(opts *app.Options) {
	// Handle global flags here
	err := opts.InitHandler()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	version, err := opts.VersionHandler()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if version != "" {
		fmt.Println(version)
		os.Exit(0)
	}
}
